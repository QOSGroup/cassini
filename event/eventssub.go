//Package event 从区块链节点监听event,
//此处只监听跨链交易event
package event

import (
	"context"
	"time"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/route"
	ctypes "github.com/QOSGroup/cassini/types"
	"github.com/tendermint/go-amino"
	pubsub "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
	ttypes "github.com/tendermint/tendermint/types"
	"strings"
)

func StartEventSubscibe(conf *config.Config) (cancel context.CancelFunc, err error) {

	for _, qsconfig := range config.DefaultQscConfig() {
		for _, nodeAddr := range strings.Split(qsconfig.NodeAddress, ",") {
			go EventSubscribe("tcp://" + nodeAddr)
		}
	}

	return
}

//EventSubscribe 从websocket服务端订阅event
//remote 服务端地址 example  "tcp://192.168.168.27:26657"
func EventSubscribe(remote string) (context.CancelFunc, error) {

	txs := make(chan interface{})

	cancel, err := SubscribeRemote(remote, "cassini", "tm.event = 'Tx'", txs)
	if err != nil {
		log.Errorf("Remote [%s] : '%s'\n", remote, err)
	}
	defer cancel()

	go func() {
		for ed := range txs {

			eventData := ed.(ttypes.EventDataTx)
			log.Infof("received event '%s'", eventData)
			cassiniEventDataTx := ctypes.CassiniEventDataTx{}

			cassiniEventDataTx.ConstructFromTags(eventData.Result.Tags)

			event := ctypes.Event{NodeAddress: remote, CassiniEventDataTx: cassiniEventDataTx}

			err := route.Event2queue(&event)

			if err != nil {
				log.Error("failed route event to message queue")
			}
		}
	}()

	return cancel, nil
}

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
