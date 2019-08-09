package commands

import (
	"github.com/spf13/cobra"
)

// NewResetCommand 创建 reset/重置（清理） 命令
func NewResetCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdReset,
		Short: "!!!WARN It's DANGER!!! Reset(cleaning up) data for relay service.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	return cmd
}
