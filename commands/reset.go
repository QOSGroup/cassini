package commands

import (
	"github.com/QOSGroup/cassini/config"
	"github.com/spf13/cobra"
)

func addResetFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&config.GetConfig().ConfigFile, "config", "./config/config.conf", "config file path")
	cmd.Flags().StringVar(&config.GetConfig().LogConfigFile, "log", "./config/log.conf", "log config file path")
}

// NewResetCommand 创建 reset/重置（清理） 命令
func NewResetCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandReset,
		Short: "!!!WARN It's DANGER!!! reset(cleaning up) data for relay service.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	addResetFlags(cmd)
	return cmd
}
