package commands

import (
	"github.com/spf13/cobra"
)

var eventNode string
var eventSubscribe string

func addEventsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&eventNode, "node", "127.0.0.1:26657", "node address")
	cmd.Flags().StringVar(&eventSubscribe, "subscribe", "tm.event='Tx' AND qcp.to='qos'", "event subscribe query")
}

// NewEventsCommand 创建 events 命令
func NewEventsCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Start web socket client and subscribe tx event",
		RunE: func(cmd *cobra.Command, args []string) error {
			mock := reconfigMock(eventNode)
			mock.Subscribe = eventSubscribe
			return commandRunner(run, isKeepRunning)
		},
	}

	addEventsFlags(cmd)
	return cmd
}
