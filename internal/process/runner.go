package process

import (
	"log"
	"os"
	"os/exec"
)

type Runner struct {
	cmd     *exec.Cmd
	command string
}

func NewRunner(command string) *Runner {
	return &Runner{
		command: command,
	}
}

func (r *Runner) Start() {
	r.cmd = exec.Command("sh", "-c", r.command)
	r.cmd.Stdout = os.Stdin
	r.cmd.Stderr = os.Stderr

	err := r.cmd.Start()
	if err != nil {
		log.Println("Failed to start:", err)
	}
}

func (r *Runner) Stop() {
	if r.cmd != nil && r.cmd.Process != nil {
		r.cmd.Process.Kill()
	}
}
