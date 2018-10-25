package main

import (
	"context"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/msgqueue"
)

// 命令行 start 命令执行方法
var starter = func(conf *config.Config) (cancel context.CancelFunc, err error) {

	log.Info("begin to start cassini")
	log.Debug("Qscs: ", len(conf.Qscs))

	for _, qsc := range conf.Qscs {
		log.Debugf("qsc: %s %s", qsc.Name, qsc.NodeAddress)
	}

	//var cancels []context.CancelFunc
	//var cancelFunc context.CancelFunc

	//启动事件监听 chain node
	//cancelFunc, err = event.StartEventSubscibe(conf)
	_, err = event.StartEventSubscibe(conf)
	//cancels = append(cancels, cancelFunc)
	if err != nil {
		return nil, err
	}

	//启动nats 消费
	err = msgqueue.StartQcpConsume(conf)
	if err != nil {
		return nil, err
	}

	//cancel = func() {
	//	for _, cancelJob := range cancels {
	//		if cancelJob != nil {
	//			cancelJob()
	//		}
	//	}
	//}

	log.Info("cassini started \n")
	return
}
