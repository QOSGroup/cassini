package commands

import (
	"github.com/spf13/cobra"
)

// NewTxCommand 创建 tx 命令
func NewTxCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "query or broadcast tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	return cmd
}
