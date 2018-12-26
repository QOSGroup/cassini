package commands

import (
	"github.com/spf13/cobra"
)

// NewVersionCommand create version command
func NewVersionCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	return cmd
}
