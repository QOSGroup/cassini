package commands

import (
	"github.com/spf13/cobra"
)

// NewStartCommand 创建 start/服务启动 命令
func NewStartCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdStart,
		Short: "Start relay service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}

	return cmd
}
