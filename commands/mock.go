package commands

import (
	"github.com/QOSGroup/cassini/config"
	"github.com/spf13/cobra"
)

func addMockFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&config.GetConfig().ConfigFile, "config", "./config/config.conf", "config file path")
	cmd.Flags().StringVar(&config.GetConfig().LogConfigFile, "log", "./config/log.conf", "log config file path")
}

// NewMockCommand 创建 mock/模拟服务 命令
func NewMockCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandMock,
		Short: "mock outer interfaces for relay service test.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	addMockFlags(cmd)
	return cmd
}
