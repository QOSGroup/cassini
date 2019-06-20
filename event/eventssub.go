//Package event 从区块链节点监听event,
//此处只监听跨链交易event
package event

import (
	"context"
	"time"

	ctypes "github.com/QOSGroup/cassini/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/rpc/client"
	tctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// SubscribeRemote subscribe events from remote
func SubscribeRemote(remote string, subscriber string, query string) (
	context.CancelFunc, <-chan tctypes.ResultEvent, error) {

	wsClient := client.NewHTTP(remote, "/websocket")

	cdc := amino.NewCodec()
	ctypes.RegisterCassiniTypesAmino(cdc)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	wsClient.Start()

	events, err := wsClient.Subscribe(ctx, subscriber, query)

	if err != nil {
		cancel()
		cancel = nil
	}

	return cancel, events, err
}
