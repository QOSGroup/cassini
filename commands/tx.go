package commands

import (
	"strings"

	"github.com/QOSGroup/cassini/config"
	"github.com/spf13/cobra"
)

var txNode string
var txSequence int64

func addTxFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&txNode, "node", "", "node address")
	cmd.Flags().Int64Var(&txSequence, "sequence", -1, "sequence")
}

// NewTxCommand 创建 tx 命令
func NewTxCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "query or broadcast tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.GetConfig()
			var mock *config.MockConfig
			if len(conf.Mocks) < 1 {
				mock = &config.MockConfig{
					RPC: &config.RPCConfig{
						ListenAddress: txNode}}
				conf.Mocks = []*config.MockConfig{mock}
			}
			if !strings.EqualFold(txNode, "") {
				if mock == nil {
					conf.Mocks = conf.Mocks[:1]
					mock = conf.Mocks[0]
				}
				mock.RPC.ListenAddress = txNode
			}
			for _, mockConf := range conf.Mocks {
				if mockConf.RPC == nil {
					mockConf.RPC = &config.RPCConfig{}
				}
				if strings.EqualFold(mockConf.RPC.ListenAddress, "") {
					mockConf.RPC.ListenAddress = DefaultNode
				}
			}
			if txSequence > -1 {
				for _, mc := range conf.Mocks {
					mc.Sequence = txSequence
				}
			}
			return commandRunner(run, isKeepRunning)
		},
	}
	addTxFlags(cmd)
	return cmd
}
