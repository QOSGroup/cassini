package commands

import (
	"github.com/spf13/cobra"
)

// NewMockCommand 创建 mock/模拟服务 命令
func NewMockCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdMock,
		Short: "Mock outer interfaces for relay service test.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	return cmd
}
