package commands

import (
	"github.com/spf13/cobra"

	"github.com/QOSGroup/cassini/config"
)

func addEventsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&config.GetConfig().EventsListen, "listen", "tcp://127.0.0.1:26657", "listen address")
	cmd.Flags().StringVar(&config.GetConfig().EventsQuery, "subscribe", "tm.event='Tx' AND qcp.to='qos'", "event subscribe query")
}

// NewEventsCommand 创建 events 命令
func NewEventsCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "start web socket client and subscribe tx event",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}

	addEventsFlags(cmd)
	return cmd
}
