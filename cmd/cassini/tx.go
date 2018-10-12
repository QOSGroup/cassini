package main

import (
	"context"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
)

// 命令行 tx 命令执行方法
var txHandler = func(conf *config.Config) (cancel context.CancelFunc, err error) {
	log.Debug("Not implemented yet!")
	return
}
