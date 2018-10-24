package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/QOSGroup/cassini/adapter"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/qbase/txs"

	tmtypes "github.com/tendermint/tendermint/types"
)

// 命令行 events 命令执行方法
var events = func(conf *config.Config) (context.CancelFunc, error) {
	cancelFunc, err := subscribe(conf.EventsListen, conf.EventsQuery)
	if err != nil {
		return nil, err
	}
	cancel := func() {
		cancelFunc()
		log.Debug("Cancel events subscribe service")
	}
	return cancel, nil
}

//subscribe 从websocket服务端订阅event
//remote 服务端地址 example  "tcp://127.0.0.1:27657"
func subscribe(remote string, query string) (context.CancelFunc, error) {
	fmt.Printf("Subscribe remote: %v, query: %v\n", remote, query)
	txsChan := make(chan interface{})
	cancel, err := event.SubscribeRemote(remote, "cassini-events", query, txsChan)
	if err != nil {
		log.Errorf("Remote [%s] : '%s'", remote, err)
		return nil, err
	}
	fmt.Printf("Subscribe successful - remote: %v, query: %v\n", remote, query)
	go func() {
		for e := range txsChan {
			et := e.(tmtypes.EventDataTx) //注：e类型断言为tmtypes.EventDataTx 类型
			var from, to string
			var seq int64
			var hash []byte
			var err error
			for _, kv := range et.Result.Tags {
				if strings.EqualFold("qcp.to", string(kv.Key)) {
					to = string(kv.Value)
				}
				if strings.EqualFold("qcp.from", string(kv.Key)) {
					from = string(kv.Value)
				}
				if strings.EqualFold("qcp.sequence", string(kv.Key)) {
					seq, err = strconv.ParseInt(string(kv.Value), 10, 64)
					if err != nil {
						log.Errorf("Get Tx event error: %v", err)
					}
				}
				if strings.EqualFold("qcp.hash", string(kv.Key)) {
					hash = kv.Value
				}
			}
			tx := &txs.TxQcp{
				BlockHeight: et.Height,
				TxIndex:     int64(et.Index),
				Sequence:    seq,
				From:        from,
				To:          to}
			fmt.Printf("Got Tx event - %v hash: %x\n", adapter.StringTx(tx), hash)

		}
	}()
	return cancel, nil
}
