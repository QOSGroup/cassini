package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/QOSGroup/cassini/adapter"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"

	tmtypes "github.com/tendermint/tendermint/types"
)

// 命令行 events 命令执行方法
var events = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	var cancels []context.CancelFunc
	var cancelFunc context.CancelFunc
	for _, mockConf := range conf.Mocks {
		cancelFunc, err = subscribe(mockConf.RPC.NodeAddress, mockConf.Subscribe)
		if err != nil {
			return
		}
		cancels = append(cancels, cancelFunc)
	}
	cancel = func() {
		for _, cancelJob := range cancels {
			if cancelJob != nil {
				cancelJob()
			}
		}
		log.Debug("Cancel events subscribe service")
	}
	return
}

//subscribe 从websocket服务端订阅event
//remote 服务端地址 example  "127.0.0.1:27657"
func subscribe(remote string, query string) (context.CancelFunc, error) {
	txsChan := make(chan interface{})
	cancel, err := event.SubscribeRemote(remote, "cassini-events", query, txsChan)
	if err != nil {
		log.Errorf("Remote %s error: %s", remote, err)
		return nil, err
	}
	log.Infof("Subscribe successful - remote: %v, subscribe: %v", remote, query)
	go func() {
		for e := range txsChan {
			et := e.(tmtypes.EventDataTx) //注：e类型断言为tmtypes.EventDataTx 类型
			var from, to string
			var seq int64
			var hash []byte
			var err error
			for _, kv := range et.Result.Tags {
				if strings.EqualFold(qcp.QcpTo, string(kv.Key)) {
					to = string(kv.Value)
				}
				if strings.EqualFold(qcp.QcpFrom, string(kv.Key)) {
					from = string(kv.Value)
				}
				if strings.EqualFold(qcp.QcpSequence, string(kv.Key)) {
					seq, err = types.BytesInt64(kv.Value)
					if err != nil {
						log.Errorf("Get Tx event error: %v", err)
					}
				}
				if strings.EqualFold(qcp.QcpHash, string(kv.Key)) {
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
