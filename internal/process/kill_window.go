//go:build windows

package process

import (
	"fmt"
	"os/exec"
)

func killProcessTree(pid int) {
	exec.Command("taskKill", "/T", "/F", "/PID", fmt.Sprint(pid)).Run()
}
