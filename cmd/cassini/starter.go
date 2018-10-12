package main

import (
	"context"

	"github.com/huangdao/cassini/config"
	"github.com/huangdao/cassini/log"
		"github.com/huangdao/cassini/event"
	)

// 命令行 start 命令执行方法
var starter = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	log.Debug("Starter")

	var cancels []context.CancelFunc
	var cancelFunc context.CancelFunc

	//启动事件监听 chain node
	cancelFunc , err = event.StartSubscibe(conf)
	cancels = append(cancels, cancelFunc)
	if err != nil {
		return nil , err
	}

	////启动nats 消费
	//err = msgqueue.StartQcpConsume(conf)
	//if err != nil {
	//	return cancelFunc , err
	//}
	//cancels = append(cancels, cancelFunc)

	cancel = func() {
		for _, cancelJob := range cancels {
			if cancelJob != nil {
				cancelJob()
			}
		}
	}
	return
}
