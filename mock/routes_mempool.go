package mock

// copy from tendermint/rpc/core/mempool.go

import (
	"context"
	"fmt"
	"time"

	"github.com/QOSGroup/cassini/log"
	"github.com/pkg/errors"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

//-----------------------------------------------------------------------------
// NOTE: tx should be signed, but this is only checked at the app level (not by Tendermint!)

// BroadcastTxAsync Returns right away, with no response
//
// ```shell
// curl 'localhost:26657/broadcast_tx_async?tx="123"'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.BroadcastTxAsync("123")
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": {
// 		"hash": "E39AAB7A537ABAA237831742DCE1117F187C3C52",
// 		"log": "",
// 		"data": "",
// 		"code": 0
// 	},
// 	"id": "",
// 	"jsonrpc": "2.0"
// }
// ```
//
// ### Query Parameters
//
// | Parameter | Type | Default | Required | Description     |
// |-----------+------+---------+----------+-----------------|
// | tx        | Tx   | nil     | true     | The transaction |
func BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	err := mempool.CheckTx(tx, nil)
	if err != nil {
		return nil, fmt.Errorf("Error broadcasting transaction: %v", err)
	}
	return &ctypes.ResultBroadcastTx{Hash: tx.Hash()}, nil
}

// BroadcastTxSync Returns with the response from CheckTx.
//
// ```shell
// curl 'localhost:26657/broadcast_tx_sync?tx="456"'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.BroadcastTxSync("456")
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"jsonrpc": "2.0",
// 	"id": "",
// 	"result": {
// 		"code": 0,
// 		"data": "",
// 		"log": "",
// 		"hash": "0D33F2F03A5234F38706E43004489E061AC40A2E"
// 	},
// 	"error": ""
// }
// ```
//
// ### Query Parameters
//
// | Parameter | Type | Default | Required | Description     |
// |-----------+------+---------+----------+-----------------|
// | tx        | Tx   | nil     | true     | The transaction |
func BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.Response, 1)
	err := mempool.CheckTx(tx, func(res *abci.Response) {
		resCh <- res
	})
	if err != nil {
		return nil, fmt.Errorf("Error broadcasting transaction: %v", err)
	}

	res := <-resCh
	r := res.GetCheckTx()
	return &ctypes.ResultBroadcastTx{
		Code: r.Code,
		Data: r.Data,
		Log:  r.Log,
		Hash: tx.Hash(),
	}, nil
}

// BroadcastTxCommit CONTRACT: only returns error if mempool.BroadcastTx errs (ie. problem with the app)
// or if we timeout waiting for tx to commit.
// If CheckTx or DeliverTx fail, no error will be returned, but the returned result
// will contain a non-OK ABCI code.
//
// ```shell
// curl 'localhost:26657/broadcast_tx_commit?tx="789"'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.BroadcastTxCommit("789")
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
// 	"error": "",
// 	"result": {
// 		"height": 26682,
// 		"hash": "75CA0F856A4DA078FC4911580360E70CEFB2EBEE",
// 		"deliver_tx": {
// 			"log": "",
// 			"data": "",
// 			"code": 0
// 		},
// 		"check_tx": {
// 			"log": "",
// 			"data": "",
// 			"code": 0
// 		}
// 	},
// 	"id": "",
// 	"jsonrpc": "2.0"
// }
// ```
//
// ### Query Parameters
//
// | Parameter | Type | Default | Required | Description     |
// |-----------+------+---------+----------+-----------------|
// | tx        | Tx   | nil     | true     | The transaction |
func BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	// subscribe to tx being committed in block
	ctx, cancel := context.WithTimeout(context.Background(), subscribeTimeout)
	defer cancel()
	deliverTxResCh := make(chan interface{})
	// q := types.EventQueryTxFor(tx)
	q := tmquery.MustParse("tm.event='Tx'")
	err := eventBus.Subscribe(ctx, "mempool", q, deliverTxResCh)
	if err != nil {
		err = errors.Wrap(err, "failed to subscribe to tx")
		log.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("Error on broadcastTxCommit: %v", err)
	}
	defer eventBus.Unsubscribe(context.Background(), "mempool", q)

	// broadcast the tx and register checktx callback
	checkTxResCh := make(chan *abci.Response, 1)
	err = mempool.CheckTx(tx, func(res *abci.Response) {
		checkTxResCh <- res
	})
	if err != nil {
		log.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("Error on broadcastTxCommit: %v", err)
	}

	// reqres := abcicli.NewReqRes(abci.ToRequestCheckTx(tx))
	// // reqres := abcicli.NewReqRes(abci.ToRequestFlush())
	// reqres.Wait()
	// checkTxResCh <- reqres.Response.GetCheckTx()
	checkTxRes := <-checkTxResCh
	checkTxR := checkTxRes.GetCheckTx()
	if checkTxR.Code != abci.CodeTypeOK {
		// CheckTx failed!
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxR,
			DeliverTx: abci.ResponseDeliverTx{},
			Hash:      tx.Hash(),
		}, nil
	}

	// reqres := abcicli.NewReqRes(abci.ToRequestCheckTx(tx))
	// checkTxR := reqres.Response.GetCheckTx()
	// checkTxR.Code = abci.CodeTypeOK
	// checkTxR := abci.ResponseCheckTx{Code: code.CodeTypeOK}

	// 构造Tx 事件（DeliverTx），模拟交易入块和发布
	tx = types.Tx("foo")
	result := abci.ResponseDeliverTx{Data: []byte("bar"), Tags: []cmn.KVPair{{Key: []byte("baz"), Value: []byte("1")}}}

	e := types.EventDataTx{types.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx,
		Result: result,
	}}

	err = eventBus.PublishEventTx(e)

	if err != nil {
		log.Error("Publish event tx error: ", err)
	}

	// Wait for the tx to be included in a block,
	// timeout after something reasonable.
	// TODO: configurable?
	timer := time.NewTimer(60 * 2 * time.Second)
	select {
	case deliverTxResMsg := <-deliverTxResCh:
		deliverTxRes := deliverTxResMsg.(types.EventDataTx)
		// The tx was included in a block.
		deliverTxR := deliverTxRes.Result
		log.Debugf("Got deliverTx event - tx: %v, response: %v", cmn.HexBytes(tx), deliverTxR)
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxR,
			DeliverTx: deliverTxR,
			Hash:      tx.Hash(),
			Height:    deliverTxRes.Height,
		}, nil
	case <-timer.C:
		log.Error("failed to include tx")
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxR,
			DeliverTx: abci.ResponseDeliverTx{},
			Hash:      tx.Hash(),
		}, fmt.Errorf("Timed out waiting for transaction to be included in a block")
	}
}

// UnconfirmedTxs Get unconfirmed transactions (maximum ?limit entries) including their number.
//
// ```shell
// curl 'localhost:26657/unconfirmed_txs'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.UnconfirmedTxs()
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
//   "error": "",
//   "result": {
//     "txs": [],
//     "n_txs": 0
//   },
//   "id": "",
//   "jsonrpc": "2.0"
// }
//
// ### Query Parameters
//
// | Parameter | Type | Default | Required | Description                          |
// |-----------+------+---------+----------+--------------------------------------|
// | limit     | int  | 30      | false    | Maximum number of entries (max: 100) |
// ```
func UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	// reuse per_page validator
	limit = validatePerPage(limit)

	txs := mempool.Reap(limit)
	return &ctypes.ResultUnconfirmedTxs{N: len(txs), Txs: txs}, nil
}

// NumUnconfirmedTxs Get number of unconfirmed transactions.
//
// ```shell
// curl 'localhost:26657/num_unconfirmed_txs'
// ```
//
// ```go
// client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
// result, err := client.UnconfirmedTxs()
// ```
//
// > The above command returns JSON structured like this:
//
// ```json
// {
//   "error": "",
//   "result": {
//     "txs": null,
//     "n_txs": 0
//   },
//   "id": "",
//   "jsonrpc": "2.0"
// }
// ```
func NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return &ctypes.ResultUnconfirmedTxs{N: mempool.Size()}, nil
}
