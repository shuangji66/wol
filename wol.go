package main

import (
	"fmt"
	"net"

	"github.com/mdlayher/wol"
)

// WakeDevice sends a magic packet to the specified MAC address.
// port can be 7 or 9.
func WakeDevice(macStr string, port int) error {
	addr, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %v", err)
	}

	// Create a new WOL client
	c, err := wol.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create WOL client: %v", err)
	}
	defer c.Close()

	// Standard WOL uses UDP broadcast to 255.255.255.255 or a specific subnet broadcast.
	// We'll use the default broadcast address.
	target := fmt.Sprintf("255.255.255.255:%d", port)
	
	err = c.Wake(target, addr)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %v", err)
	}

	return nil
}
