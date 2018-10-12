package mock

import (
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

// copy from tendermint/state/service.go

// MempoolMocker mock 实现，用于事件订阅和RPC 服务模拟
type MempoolMocker struct {
}

// Lock 接口方法实现
func (m MempoolMocker) Lock() {}

// Unlock 接口方法实现
func (m MempoolMocker) Unlock() {}

// Size 接口方法实现
func (m MempoolMocker) Size() int { return 0 }

// CheckTx 接口方法实现
func (m MempoolMocker) CheckTx(tx types.Tx, cb func(*abci.Response)) error {
	// reqres := abcicli.NewReqRes(abci.ToRequestCheckTx(tx))
	// reqres.SetCallback(cb)
	// reqres.Wait()
	// req := abci.ToRequestCheckTx(tx)
	// res, err := cli.client.CheckTx(context.Background(), req.GetCheckTx())
	// if err != nil {
	// 	cli.StopForError(err)
	// }
	// cb(&types.Response{Value: &types.Response_CheckTx{res}})
	// cb(&abci.Response{Value: &abci.Response_CheckTx{&abci.RequestCheckTx{Tx: tx}}})
	// checkTx := req.GetCheckTx()
	out := new(abci.ResponseCheckTx)
	out.Code = abci.CodeTypeOK
	out.Data = []byte("foo")
	out.Tags = []cmn.KVPair{{Key: []byte("baz"), Value: []byte("1")}}

	cb(&abci.Response{Value: &abci.Response_CheckTx{CheckTx: out}})
	return nil
}

// Reap 接口方法实现
func (m MempoolMocker) Reap(n int) types.Txs { return types.Txs{} }

// Update 接口方法实现
func (m MempoolMocker) Update(height int64, txs types.Txs) error { return nil }

// Flush 接口方法实现
func (m MempoolMocker) Flush() {}

// FlushAppConn 接口方法实现
func (m MempoolMocker) FlushAppConn() error { return nil }

// TxsAvailable 接口方法实现
func (m MempoolMocker) TxsAvailable() <-chan struct{} { return make(chan struct{}) }

// EnableTxsAvailable 接口方法实现
func (m MempoolMocker) EnableTxsAvailable() {}
