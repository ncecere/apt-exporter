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
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	dir := filepath.Dir(execPath)

	// Build the path to the new executable
	newExePath := filepath.Join(dir, "cmd", "apt-exporter")

	// Check if the new executable exists as a compiled binary
	if _, err := os.Stat(newExePath); err == nil {
		// Forward all arguments to the new executable
		cmd := exec.Command(newExePath, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Run the new executable
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to execute new binary: %v", err)
		}
		return
	}

	// If the compiled binary doesn't exist, try running the Go file directly
	newGoFilePath := filepath.Join(dir, "cmd", "apt-exporter", "main.go")
	if _, err := os.Stat(newGoFilePath); os.IsNotExist(err) {
		log.Fatalf("New executable not found at %s or %s: %v", newExePath, newGoFilePath, err)
	}

	// Forward all arguments to the new executable using go run
	cmd := exec.Command("go", append([]string{"run", newGoFilePath}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the new executable
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute new binary: %v", err)
	}
}
