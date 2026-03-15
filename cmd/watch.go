package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cli/internal/process"
	"github.com/spf13/cli/internal/watcher"
	"github.com/spf13/cobra"
)

var command string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch files and restart process on change",
	Run: func(cmd *cobra.Command, args []string) {
		timer := time.NewTimer(0)
		timer.Stop()
		events := make(chan struct{}, 1)
		runner := process.NewRunner(command)
		runner.Start()

		go watcher.StartWatching(events)
		fmt.Println("Watching files...")
		fmt.Println("Command:", command)

		for {
			select {
			case <-events:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(time.Millisecond * 300)

			case <-timer.C:
				fmt.Println("Restarting process...")
				runner.Stop()
				runner.Start()

			}

		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run")
	watchCmd.MarkFlagRequired("command")
}
