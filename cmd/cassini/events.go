package main

import (
	"context"
	"fmt"
	"os"

	"github.com/QOSGroup/cassini/adapter/ports"
	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/QOSGroup/qbase/txs"
)

// 命令行 events 命令执行方法
var events = func() (cancel context.CancelFunc, err error) {
	conf := config.GetConfig()
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
	ip, port, err := ports.ParseNodeAddress(remote)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf := &ports.AdapterConfig{
		ChainName: "cassini-events",
		ChainType: "qos",
		IP:        ip,
		Port:      port,
		Query:     query}
	conf.Listener = func(event *types.Event, adapter ports.Adapter) {
		handle(event)
	}
	ports.RegisterAdapter(conf)
	log.Infof("Subscribe successful - remote: %v, subscribe: %v", remote, query)

	cancel := func() {
	}
	return cancel, nil
}

func handle(event *types.Event) {
	// et := e.(tmtypes.EventDataTx) //注：e类型断言为tmtypes.EventDataTx 类型
	// ca := types.CassiniEventDataTx{}
	// err := ca.ConstructFromTags(et.Result.Tags)
	ca := event.CassiniEventDataTx
	tx := &txs.TxQcp{
		BlockHeight: event.Height,
		// TxIndex:     int64(et.Index),
		Sequence: ca.Sequence,
		From:     ca.From,
		To:       ca.To}
	log.Debugf("Got Tx event - %v hash: %x\n",
		cmn.StringTx(tx), ca.HashBytes)
}
