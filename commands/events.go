package commands

import (
	"fmt"
	"strings"

	"github.com/QOSGroup/cassini/config"
	"github.com/spf13/cobra"
)

var eventNode string
var eventSubscribe string

func addEventsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&eventNode, "node", "", "node address")
	cmd.Flags().StringVar(&eventSubscribe, "subscribe", "tm.event='Tx' AND qcp.to='qos'", "event subscribe query")
}

// NewEventsCommand 创建 events 命令
func NewEventsCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "start web socket client and subscribe tx event",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.GetConfig()
			var mock *config.MockConfig
			if len(conf.Mocks) < 1 {
				mock = &config.MockConfig{
					RPC: &config.RPCConfig{
						ListenAddress: eventNode}}
				conf.Mocks = []*config.MockConfig{mock}
			}
			if !strings.EqualFold(eventNode, "") {
				if mock == nil {
					conf.Mocks = conf.Mocks[:1]
					mock = conf.Mocks[0]
				}
				mock.RPC.ListenAddress = eventNode
			}
			for _, mockConf := range conf.Mocks {
				if mockConf.RPC == nil {
					mockConf.RPC = &config.RPCConfig{}
				}
				if strings.EqualFold(mockConf.RPC.ListenAddress, "") {
					mockConf.RPC.ListenAddress = DefaultNode
				}
			}
			if !strings.EqualFold(eventSubscribe, "") {
				for _, mc := range conf.Mocks {
					mc.Subscribe = eventSubscribe
				}
			}
			fmt.Println("events RPC: ", mock.RPC.ListenAddress)
			return commandRunner(run, isKeepRunning)
		},
	}

	addEventsFlags(cmd)
	return cmd
}
