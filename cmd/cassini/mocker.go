package main

import (
	"context"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/mock"
)

// 命令行 mock 命令执行方法
var mocker = func() (cancel context.CancelFunc, err error) {
	log.Debug("Load mocks: ", len(config.GetConfig().Mocks))
	var cancels []context.CancelFunc
	var cancelFunc context.CancelFunc
	for _, mockConf := range config.GetConfig().Mocks {
		cancelFunc, err = mock.StartMock(*mockConf)
		if err != nil {
			return
		}
		cancels = append(cancels, cancelFunc)

		//time.Sleep(1100 * time.Millisecond)
	}
	//cancelFunc, err = starter(conf)
	//if err != nil {
	//	return
	//}
	// cancels = append(cancels, cancelFunc)
	cancel = func() {
		for _, cancelJob := range cancels {
			if cancelJob != nil {
				cancelJob()
			}
		}
	}
	return
}
