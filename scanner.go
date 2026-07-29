package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ScanLAN scans the local network for active devices.
// 1. First sends UDP broadcast to populate ARP cache
// 2. Then reads ARP table to get device list
func ScanLAN() ([]Device, error) {
	log.Println("Starting LAN scan...")

	// Step 1: 发送 UDP 广播/扫描，让更多设备进入 ARP 缓存
	populateARPCache()

	// Step 2: 等待一下让 ARP 缓存更新
	time.Sleep(2 * time.Second)

	// Step 3: 读取 ARP 表
	devices, err := scanFromProcARPFile()
	if err == nil && len(devices) > 0 {
		log.Printf("Found %d devices from /proc/net/arp", len(devices))
		return devices, nil
	}
	if err != nil {
		log.Printf("/proc/net/arp failed: %v", err)
	}

	// Fallback: ip neigh
	devices, err = scanFromIPNeigh()
	if err == nil && len(devices) > 0 {
		log.Printf("Found %d devices from ip neigh", len(devices))
		return devices, nil
	}

	// Fallback: arp -an
	devices, err = scanFromArpCommand()
	if err == nil && len(devices) > 0 {
		log.Printf("Found %d devices from arp", len(devices))
		return devices, nil
	}

	log.Println("All scan methods returned no devices")
	return []Device{}, nil
}

// populateARPCache 通过向子网内所有 IP 发送 UDP 包来填充 ARP 缓存
// 不需要 root 权限，普通用户就能发 UDP
func populateARPCache() {
	localIPs, err := getLocalIPs()
	if err != nil {
		log.Printf("Failed to get local IPs: %v", err)
		return
	}

	log.Printf("Local IPs: %v", localIPs)

	var wg sync.WaitGroup
	for _, localIP := range localIPs {
		// 获取子网前缀，例如 192.168.1.
		prefix := localIP[:strings.LastIndex(localIP, ".")+1]
		log.Printf("Scanning subnet: %s0/24", prefix)

		// 向子网广播地址发 UDP 包
		broadcastAddr := prefix + "255"
		sendUDP(broadcastAddr, 9)

		// 向每个 IP 发送 UDP 包（并行）
		for i := 1; i < 255; i++ {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				sendUDP(ip, 9) // port 9 = discard protocol
			}(fmt.Sprintf("%s%d", prefix, i))
		}
	}

	wg.Wait()
	log.Println("UDP broadcast sweep completed")
}

// sendUDP 发送一个小 UDP 包到指定地址，触发 ARP 解析
func sendUDP(ip string, port int) {
	conn, err := net.DialTimeout("udp4", fmt.Sprintf("%s:%d", ip, port), 200*time.Millisecond)
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(200 * time.Millisecond))
	conn.Write([]byte{0})
	conn.Close()
}

func getLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

// scanFromProcARPFile reads /proc/net/arp using Go file I/O (no exec needed)
func scanFromProcARPFile() ([]Device, error) {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("cannot open /proc/net/arp: %v", err)
	}
	defer file.Close()

	var devices []Device
	scanner := bufio.NewScanner(file)

	// Skip header
	if scanner.Scan() {
		log.Printf("ARP header: %s", scanner.Text())
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		ip := fields[0]
		mac := fields[3]

		if mac == "00:00:00:00:00:00" {
			continue
		}

		devices = append(devices, Device{
			IP:       ip,
			MAC:      mac,
			Hostname: resolveHostname(ip),
			Online:   true,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/net/arp: %v", err)
	}

	return devices, nil
}

// scanFromIPNeigh tries common paths for the `ip` command
func scanFromIPNeigh() ([]Device, error) {
	paths := []string{"/sbin/ip", "/usr/sbin/ip", "/bin/ip", "/usr/bin/ip", "ip"}

	var output []byte
	var err error
	for _, p := range paths {
		output, err = exec.Command(p, "neigh").Output()
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	var devices []Device
	macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}`)

	sc := bufio.NewScanner(strings.NewReader(string(output)))
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[len(fields)-1] == "FAILED" {
			continue
		}
		ip := fields[0]
		mac := macRegex.FindString(line)
		if mac == "" {
			continue
		}
		devices = append(devices, Device{
			IP:       ip,
			MAC:      mac,
			Hostname: resolveHostname(ip),
			Online:   true,
		})
	}
	return devices, nil
}

// scanFromArpCommand tries common paths for the `arp` command
func scanFromArpCommand() ([]Device, error) {
	paths := []string{"/sbin/arp", "/usr/sbin/arp", "/bin/arp", "/usr/bin/arp", "arp"}

	var output []byte
	var err error
	for _, p := range paths {
		output, err = exec.Command(p, "-an").Output()
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	var devices []Device
	ipRegex := regexp.MustCompile(`\(([\d\.]+)\)`)
	macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}[:-]){5}([0-9a-fA-F]{2})`)

	sc := bufio.NewScanner(strings.NewReader(string(output)))
	for sc.Scan() {
		line := sc.Text()
		ipMatch := ipRegex.FindStringSubmatch(line)
		macMatch := macRegex.FindString(line)
		if len(ipMatch) > 1 && macMatch != "" {
			devices = append(devices, Device{
				IP:       ipMatch[1],
				MAC:      macMatch,
				Hostname: resolveHostname(ipMatch[1]),
				Online:   true,
			})
		}
	}
	return devices, nil
}

func resolveHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], ".")
	}
	return ""
}
