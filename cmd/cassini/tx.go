package main

import (
	"context"

	"github.com/huangdao/cassini/config"
	"github.com/huangdao/cassini/log"
)

// 命令行 tx 命令执行方法
var txHandler = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	log.Debug("Not implemented yet!")
	return
}
