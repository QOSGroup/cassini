package main

import (
	"context"
	"fmt"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/restclient"
)

// 命令行 tx 命令执行方法
var txHandler = func(conf *config.Config) (context.CancelFunc, error) {
	log.Debug("tx")
	client := restclient.NewRestClient("tcp://127.0.0.1:26657")
	tx, err := client.GetTxQcp("qsc", 0)
	if err != nil {
		log.Error("Get TxQcp error: ", err)
	}
	if tx != nil {
		fmt.Println("Get TxQcp from: ", tx.From, " to: ", tx.To)
	}
	return nil, nil
}
