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
	tx, err := client.GetTxQcp("qstar", 1)
	if err != nil {
		log.Error("Get TxQcp error: ", err)
	}
	if tx != nil {
		fmt.Printf("Get TxQcp chain: %s to: %s\n", tx.From, tx.To)
	}
	var seq int64
	seq, err = client.GetSequence("qstar", "out")
	if err != nil {
		log.Error("Get sequence error: ", err)
	}
	fmt.Println("Get sequence: ", seq)
	return nil, nil
}
