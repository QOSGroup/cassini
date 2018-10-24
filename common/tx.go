package common

import (
	"fmt"

	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Transform 将交易转换为交易事件
func Transform(tx *txs.TxQcp) (*tmtypes.EventDataTx, error) {
	hash := HashTx(tx)
	result := abcitypes.ResponseDeliverTx{
		Data: []byte("mock"),
		Tags: []cmn.KVPair{
			{Key: []byte(qcp.QcpTo), Value: []byte(tx.To)},
			{Key: []byte(qcp.QcpFrom), Value: []byte(tx.From)},
			{Key: []byte(qcp.QcpSequence), Value: []byte(fmt.Sprintf("%v", tx.Sequence))},
			{Key: []byte(qcp.QcpHash), Value: hash},
		}}
	return &tmtypes.EventDataTx{TxResult: tmtypes.TxResult{
		Height: tx.BlockHeight,
		Index:  uint32(tx.TxIndex),
		Tx:     tx.GetSigData(),
		Result: result,
	}}, nil
}
