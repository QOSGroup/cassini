package commands

import (
	"github.com/spf13/cobra"
)

var txNode string
var txSequence int64

func addTxFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&txNode, "node", "127.0.0.1:26657", "node address")
	cmd.Flags().Int64Var(&txSequence, "sequence", -1, "sequence")
}

// NewTxCommand 创建 tx 命令
func NewTxCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "query or broadcast tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			mock := reconfigMock(txNode)
			if txSequence > -1 {
				mock.Sequence = txSequence
			}
			return commandRunner(run, isKeepRunning)
		},
	}
	addTxFlags(cmd)
	return cmd
}
