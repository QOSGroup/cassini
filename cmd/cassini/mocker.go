package main

import (
	"context"

	"github.com/huangdao/cassini/config"
	"github.com/huangdao/cassini/log"
	"github.com/huangdao/cassini/mock"
)

// 命令行 mock 命令执行方法
var mocker = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	log.Debug("Load mocks: ", len(config.GetConfig().Mocks))
	var cancels []context.CancelFunc
	var cancelFunc context.CancelFunc
	for _, mockConf := range config.GetConfig().Mocks {
		cancelFunc, err = mock.StartQscMock(&mockConf)
		if err != nil {
			return
		}
		cancels = append(cancels, cancelFunc)
	}
	//cancelFunc, err = starter(conf)
	//if err != nil {
	//	return
	//}
	cancels = append(cancels, cancelFunc)
	cancel = func() {
		for _, cancelJob := range cancels {
			if cancelJob != nil {
				cancelJob()
			}
		}
	}
	return
}
