// Package main provides a backward compatibility wrapper for the apt-exporter.
// This file is kept for backward compatibility and will forward execution to the new location.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("WARNING: This entry point is deprecated. Please use the new binary at cmd/apt-exporter/main.go")
	fmt.Println("Forwarding execution to the new location...")

	// Get the directory of the current executable
	dir := filepath.Dir(os.Args[0])

	// Build the path to the new executable
	newExePath := filepath.Join(dir, "cmd", "apt-exporter", "main.go")

	// Check if the new executable exists
	if _, err := os.Stat(newExePath); os.IsNotExist(err) {
		log.Fatalf("New executable not found at %s: %v", newExePath, err)
	}

	// Forward all arguments to the new executable
	cmd := exec.Command("go", append([]string{"run", newExePath}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the new executable
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute new binary: %v", err)
	}
}
