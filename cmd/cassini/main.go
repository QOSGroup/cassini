// Package main 跨链中继服务程序
//
// 包括服务配置，启动服务(start)、模拟运行服务(mock)以及交易事件监听(events)
package main

import (
	"os"
	"path/filepath"

	_ "github.com/QOSGroup/cassini/adapter/ports/ethereum"
	_ "github.com/QOSGroup/cassini/adapter/ports/fabric"
	"github.com/QOSGroup/cassini/commands"
	"github.com/QOSGroup/cassini/log"
)

// 链中继服务主程序，包括正常服务启动和mock 测试运行两种模式
// 细节参数配置帮助信息请运行帮助命令查看：cassini help
func main() {
	defer log.Flush()

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(starter, true),
		commands.NewEventsCommand(events, true),
		commands.NewMockCommand(mocker, true),
		commands.NewResetCommand(resetHandler, false),
		commands.NewTxCommand(txHandler, false),
		commands.NewVersionCommand(versioner, false))

	defaultHome := os.ExpandEnv("$HOME/.cassini")
	defaultConfig := filepath.Join(defaultHome, "config/cassini.yml")
	defaultLog := filepath.Join(defaultHome, "config/log.conf")

	root.PersistentFlags().String(commands.FlagHome,
		defaultHome, "Directory for config and data")
	root.PersistentFlags().String(commands.FlagConfig,
		defaultConfig, "Config file path")
	root.PersistentFlags().String(commands.FlagLog,
		defaultLog, "Log config file path")

	if err := root.Execute(); err != nil {
		log.Error("Exit by error: ", err)
	}
}
