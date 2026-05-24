package cli

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func ensureDaemon(socketPath string) error {
	if isDaemonAlive(socketPath) {
		return nil
	}

	// Stale socket file from a dirty exit prevents net.Listen on the same path.
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	fmt.Println("Bufo daemon inicializando.")
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exe, "daemon", "serve")
	if err := cmd.Start(); err != nil {
		return err
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if isDaemonAlive(socketPath) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("daemon não iniciou a tempo")
}

func isDaemonAlive(socketPath string) bool {
	conn, err := net.DialTimeout("unix", socketPath, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
