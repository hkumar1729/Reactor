//go:build linux || darwin

package process

import (
	"fmt"
	"log"
	"syscall"
)

func killProcessTree(pid int) {
	err := syscall.Kill(-pid, syscall.SIGKILL)
	if err != nil {
		log.Fatalf("failed to kill process %d: %v\n", pid, err)
	}

	fmt.Printf("Process %d killed with SIGKILL\n", pid)
}
