//Package event 从区块链节点监听event,
//此处只监听跨链交易event
package event

import (
	"context"
	"time"

	"errors"
	"fmt"
	"strings"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/route"
	ctypes "github.com/QOSGroup/cassini/types"
	"github.com/tendermint/go-amino"
	pubsub "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
	ttypes "github.com/tendermint/tendermint/types"
)

// StartEventSubscibe Start event listening
func StartEventSubscibe(conf *config.Config) (cancel context.CancelFunc, err error) {

	var subEventFrom string
	es := make(chan error, 1024) //TODO 1024根据节点数需要修改

	// for _, qsconfig := range config.DefaultQscConfig() {
	for _, qsconfig := range conf.Qscs {
		for _, nodeAddr := range strings.Split(qsconfig.NodeAddress, ",") {
			go EventsSubscribe(conf, "tcp://"+nodeAddr, es)
			subEventFrom += fmt.Sprintf("[%s] ", nodeAddr)

		}
	}

	if len(es) > 0 {
		return nil, errors.New("subscibe events failed")
	}

	log.Infof("subscibed events from %s", subEventFrom)

	return
}

// EventsSubscribe 从websocket服务端订阅event
//remote 服务端地址 example  "tcp://192.168.168.27:26657"
func EventsSubscribe(conf *config.Config, remote string, e chan<- error) context.CancelFunc {
	log.Debug("Event subscribe remote: ", remote)

	txs := make(chan interface{})
	//TODO query 条件?? "tm.event = 'Tx' AND qcp.to != '' AND qcp.sequence > 0"
	cancel, err := SubscribeRemote(remote,
		"cassini", "tm.event = 'Tx'  AND qcp.sequence > 0", txs)
	if err != nil {
		e <- err
		log.Errorf("Remote [%s] : '%s'", remote, err)
	}
	//defer cancel() //TODO  panic

	go func() {
		for ed := range txs {

			eventData := ed.(ttypes.EventDataTx)
			log.Debugf("Received event from[%s],'%s'", remote, eventData)

			cassiniEventDataTx := ctypes.CassiniEventDataTx{}
			cassiniEventDataTx.Height = eventData.Height
			cassiniEventDataTx.ConstructFromTags(eventData.Result.Tags)

			event := ctypes.Event{NodeAddress: remote, CassiniEventDataTx: cassiniEventDataTx}

			_, err := route.Event2queue(conf, &event)

			if err != nil {
				log.Errorf("failed route event to message queue,%s", err.Error())
			}
		}
	}()

	return cancel
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
