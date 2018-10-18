package main

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"

	tmtypes "github.com/tendermint/tendermint/types"
)

// 命令行 events 命令执行方法
var events = func(conf *config.Config) (context.CancelFunc, error) {
	cancelFunc, err := Subscribe(conf.EventsListen, conf.EventsQuery)
	if err != nil {
		return nil, err
	}
	cancel := func() {
		cancelFunc()
		log.Debug("Cancel events subscribe service")
	}
	return cancel, nil
}

//Subscribe 从websocket服务端订阅event
//remote 服务端地址 example  "tcp://127.0.0.1:26657"
func Subscribe(remote string, query string) (context.CancelFunc, error) {
	fmt.Printf("Subscribe remote: %v, query: %v\n", remote, query)
	txs := make(chan interface{})
	cancel, err := event.SubscribeRemote(remote, "cassini-events", query, txs)
	if err != nil {
		log.Errorf("Remote [%s] : '%s'", remote, err)
		return nil, err
	}
	fmt.Printf("Subscribe successful - remote: %v, query: %v\n", remote, query)
	go func() {
		for e := range txs {
			fmt.Println("Got Tx event - ", e.(tmtypes.EventDataTx)) //注：e类型断言为types.CassiniEventDataTx 类型
			for _, tto := range e.(tmtypes.EventDataTx).Result.Tags {
				kv := tto //interface{}(tto).(common.KVPair)
				fmt.Println(string(kv.Key), string(kv.Value))
			}
		}
	}()
	return cancel, nil
}
