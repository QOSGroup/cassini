package commands

import (
	"github.com/spf13/cobra"
)

func addMockFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&conf, "config", "./config/mock_config.conf", "config file path　for mock mode")
	cmd.Flags().StringVar(&logConf, "log", "./config/mock_log.conf", "log config file path for mock mode")
}

// NewMockCommand 创建 mock/模拟服务 命令
func NewMockCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mock",
		Short: "mock outer interfaces for relay service test.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}

	addMockFlags(cmd)
	return cmd
}
