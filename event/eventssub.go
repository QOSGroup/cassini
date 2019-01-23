//Package event 从区块链节点监听event,
//此处只监听跨链交易event
package event

import (
	"context"
	"time"

	ctypes "github.com/QOSGroup/cassini/types"
	"github.com/tendermint/go-amino"
	pubsub "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
)

// SubscribeRemote 订阅接口，暴露检测点以便于测试
func SubscribeRemote(remote string, subscriber string, query string, txs chan<- interface{}) (context.CancelFunc, error) {

	wsClient := client.NewHTTP(remote, "/websocket")

	cdc := amino.NewCodec()
	ctypes.RegisterCassiniTypesAmino(cdc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	//query := query.MustParse("tm.event = 'Tx' AND tx.height = 3")
	q := pubsub.MustParse(query)
	wsClient.Start()

	err := wsClient.Subscribe(ctx, subscriber, q, txs) //注：不仅订阅 还完成了event的amino解码 在httpclient.go 函数eventListener

	if err != nil {
		cancel()
		cancel = nil
	}

	return cancel, err
}
