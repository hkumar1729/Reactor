package cmd

import (
	"fmt"

	"github.com/spf13/cli/internal/process"
	"github.com/spf13/cli/internal/watcher"
	"github.com/spf13/cobra"
)

var command string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch files and restart process on change",
	Run: func(cmd *cobra.Command, args []string) {
		events := make(chan struct{})
		runner := process.NewRunner(command)
		runner.Start()

		go watcher.StartWatching(events)
		fmt.Println("Watching files...")
		fmt.Println("Command:", command)

		for {
			<-events
			runner.Stop()
			runner.Start()
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run")
	watchCmd.MarkFlagRequired("command")
}
