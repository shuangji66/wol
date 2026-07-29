package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

//go:embed dist/*
var distFS embed.FS

type Device struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Online   bool   `json:"online"`
	Group    string `json:"group"`
	Schedule string `json:"schedule"`
}

type Config struct {
	Password string
	AutoScan bool
}

type App struct {
	sync.Mutex
	Devices  []Device
	Config   Config
	dataDir  string
}

func (a *App) devicesFilePath() string {
	return filepath.Join(a.dataDir, "devices.json")
}

func (a *App) loadDevices() {
	file, err := os.ReadFile(a.devicesFilePath())
	if err == nil {
		json.Unmarshal(file, &a.Devices)
	}
}

func (a *App) saveDevices() {
	data, _ := json.Marshal(a.Devices)
	path := a.devicesFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Printf("ERROR creating data directory: %v", err)
		return
	}
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		log.Printf("ERROR writing devices.json to %s: %v", path, err)
	} else {
		log.Printf("Successfully saved %d devices to %s", len(a.Devices), path)
	}
}

func (a *App) startScheduler() {
	ticker := time.NewTicker(20 * time.Second)
	go a.updateDeviceStatus()
	for range ticker.C {
		now := time.Now().Format("15:04")
		a.Lock()
		for _, d := range a.Devices {
			if d.Schedule == now {
				log.Printf("Scheduled wakeup for %s (%s)", d.Hostname, d.MAC)
				WakeDevice(d.MAC, 9)
			}
		}
		a.Unlock()
		a.updateDeviceStatus()
	}
}

func (a *App) updateDeviceStatus() {
	populateARPCache()
	time.Sleep(2 * time.Second)
	activeDevices, err := scanFromProcARPFile()
	if err != nil || len(activeDevices) == 0 {
		activeDevices, _ = scanFromIPNeigh()
	}
	if len(activeDevices) == 0 {
		activeDevices, _ = scanFromArpCommand()
	}
	activeMap := make(map[string]bool)
	for _, dev := range activeDevices {
		activeMap[strings.ToLower(dev.MAC)] = true
	}
	a.Lock()
	changed := false
	for i, dev := range a.Devices {
		isOnline := activeMap[strings.ToLower(dev.MAC)]
		if dev.Online != isOnline {
			a.Devices[i].Online = isOnline
			changed = true
		}
	}
	if changed {
		a.saveDevices()
		log.Println("Device online status updated")
	}
	a.Unlock()
}

func main() {
	dataDir := flag.String("d", ".", "Directory to store devices.json")
	flag.Parse()

	socketPath := os.Getenv("WOL_SOCKET")
	if socketPath == "" {
		socketPath = "./wol.sock"
	}
	baseURL := os.Getenv("WOL_BASE_URL")
	if baseURL != "" && !strings.HasPrefix(baseURL, "/") {
		baseURL = "/" + baseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("=== WOL Server starting ===")
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("Socket: %s, BaseURL: %s", socketPath, baseURL)

	app := &App{
		dataDir: *dataDir,
	}
	app.loadDevices()
	log.Printf("Loaded %d devices", len(app.Devices))
	app.Config.Password = os.Getenv("WOL_PASSWORD")
	app.Config.AutoScan = os.Getenv("WOL_AUTO_SCAN") == "true"
	go app.startScheduler()

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[API] %s /api/devices", r.Method)
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			app.Lock()
			devs := app.Devices
			if devs == nil {
				devs = []Device{}
			}
			json.NewEncoder(w).Encode(devs)
			app.Unlock()
		} else if r.Method == "POST" {
			var req struct {
				Device
				OldMAC string `json:"oldMac"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			app.Lock()
			found := false
			searchMac := req.MAC
			if req.OldMAC != "" {
				searchMac = req.OldMAC
			}
			for i, dev := range app.Devices {
				if strings.EqualFold(strings.TrimSpace(dev.MAC), strings.TrimSpace(searchMac)) {
					app.Devices[i] = req.Device
					found = true
					break
				}
			}
			if !found {
				app.Devices = append(app.Devices, req.Device)
			}
			app.saveDevices()
			app.Unlock()
			w.WriteHeader(200)
		} else if r.Method == "DELETE" {
			mac := r.URL.Query().Get("mac")
			app.Lock()
			for i, dev := range app.Devices {
				if dev.MAC == mac {
					app.Devices = append(app.Devices[:i], app.Devices[i+1:]...)
					break
				}
			}
			app.saveDevices()
			app.Unlock()
			w.WriteHeader(200)
		}
	})

	mux.HandleFunc("/api/scan", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[API] POST /api/scan - starting scan")
		w.Header().Set("Content-Type", "application/json")
		devices, err := ScanLAN()
		if err != nil {
			log.Printf("[API] Scan failed: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("[API] Scan found %d devices", len(devices))
		json.NewEncoder(w).Encode(devices)
	})

	mux.HandleFunc("/api/shutdown", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			IP   string `json:"ip"`
			User string `json:"user"`
			Pass string `json:"pass"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		err := ShutdownDevice(req.IP, req.User, req.Pass)
		if err != nil {
			http.Error(w, fmt.Sprintf("Shutdown failed: %v", err), 500)
			return
		}
		w.WriteHeader(200)
	})

	mux.HandleFunc("/api/wake", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			MAC   string `json:"mac"`
			Group string `json:"group"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		app.Lock()
		defer app.Unlock()
		if req.Group != "" {
			for _, d := range app.Devices {
				if d.Group == req.Group {
					WakeDevice(d.MAC, 9)
				}
			}
		} else {
			WakeDevice(req.MAC, 9)
		}
		w.WriteHeader(200)
	})

	// Static file server (embedded dist)
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	// Outer handler for Base URL stripping and redirect
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if baseURL != "" && path == baseURL && !strings.HasSuffix(path, "/") {
			newPath := path + "/"
			if r.URL.RawQuery != "" {
				newPath += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, newPath, http.StatusMovedPermanently)
			return
		}

		if baseURL != "" {
			if !strings.HasPrefix(path, baseURL) {
				http.NotFound(w, r)
				return
			}
			r.URL.Path = strings.TrimPrefix(path, baseURL)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
		}

		mux.ServeHTTP(w, r)
	})

	// Remove old socket if exists
	if err := os.RemoveAll(socketPath); err != nil {
		log.Printf("Warning: could not remove existing socket: %v", err)
	}

	// Create listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(socketPath, 0666); err != nil {
		log.Printf("Warning: could not chmod socket: %v", err)
	}

	// Ensure socket is removed when program exits
	defer func() {
		if err := os.Remove(socketPath); err != nil {
			log.Printf("Warning: failed to remove socket on exit: %v", err)
		} else {
			log.Printf("Socket removed: %s", socketPath)
		}
	}()

	// Signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Listening on unix socket: %s", socketPath)
		if err := http.Serve(listener, handler); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP serve error: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigCh
	log.Println("Received shutdown signal, stopping server...")
	listener.Close()
	// http.Serve will return after listener.Close(), then defer will remove the socket
}