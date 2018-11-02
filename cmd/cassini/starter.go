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

	log.Tracef("Qscs: %d", len(conf.Qscs))
	for _, qsc := range conf.Qscs {
		log.Tracef("qsc: %s %s", qsc.Name, qsc.NodeAddress)
	}

	//启动事件监听 chain node
	_, err = event.StartEventSubscibe(conf)
	if err != nil {
		log.Errorf("Cassini start error: %v", err)
	}

	//启动nats 消费
	err = msgqueue.StartQcpConsume(conf)
	if err != nil {
		return nil, err
	}

	log.Info("cassini started \n")
	return
}
