// Package main 跨链中继服务程序
//
// 包括服务配置，启动服务(start)、模拟运行服务(mock)以及交易事件监听(events)
package main

import (
	"context"
	"fmt"
	"github.com/QOSGroup/cassini/commands"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/version"
)

// 链中继服务主程序，包括正常服务启动和mock 测试运行两种模式
// 细节参数配置帮助信息请运行帮助命令查看：cassini help
func main() {
	defer log.Flush()

	root := commands.NewRootCommand()
	root.AddCommand(
		commands.NewStartCommand(starter, true),
		commands.NewMockCommand(mocker, true),
		commands.NewEventsCommand(events, true),
		commands.NewTxCommand(txHandler, false),
		commands.NewVersionCommand(versioner, false))

	if err := root.Execute(); err != nil {
		log.Error("Exit by error: ", err)
	}
	log.Debug("Ok.")
}

var versioner = func(conf *config.Config) (context.CancelFunc, error) {

	fmt.Println(version.Version)
	return nil, nil
}
