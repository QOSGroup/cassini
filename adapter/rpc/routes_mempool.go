package rpc

// copy from tendermint/rpc/core/mempool.go

import (
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

//-----------------------------------------------------------------------------
// NOTE: tx should be signed, but this is only checked at the app level (not by Tendermint!)

// BroadcastTxSync 广播交易。
func (s RequestHandler) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return &ctypes.ResultBroadcastTx{
		Code: abci.CodeTypeOK,
		Hash: tx.Hash(),
	}, nil
}
