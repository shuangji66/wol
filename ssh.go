package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

// RemoteCommand executes a command via SSH on the target device.
func RemoteCommand(ip, user, password, cmd string) error {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	return nil
}

// ShutdownDevice tries to shutdown or sleep the device.
// Typical commands:
// Windows: "shutdown /s /t 0"
// Linux: "systemctl poweroff" or "shutdown -h now"
func ShutdownDevice(ip, user, password string) error {
	// Let's try multiple common commands or allow user to configure
	return RemoteCommand(ip, user, password, "sudo systemctl poweroff || sudo shutdown -h now || shutdown /s /t 0")
}
