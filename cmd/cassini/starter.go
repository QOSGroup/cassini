package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/QOSGroup/cassini/concurrency"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/msgqueue"
)

// 命令行 start 命令执行方法
var starter = func(conf *config.Config) (cancel context.CancelFunc, err error) {

	log.Info("Starting cassini...")

	var w sync.WaitGroup
	w.Add(1)
	go func() {
		etcd, e := concurrency.StartEmbedEtcd(conf)
		if e != nil {
			log.Error("Etcd server start error: ", e)
			log.Flush()
			os.Exit(1)
		}
		w.Done()
		if etcd == nil {
			return
		}
		defer etcd.Close()
		select {
		case <-etcd.Server.ReadyNotify():
			log.Info("Etcd server is ready!")
		case <-time.After(60 * time.Second):
			etcd.Server.Stop() // trigger a shutdown
			log.Info("Etcd server took too long to start!")
		}
		e = <-etcd.Err()
		log.Error("Etcd running error: ", e)
	}()

	log.Tracef("Qscs: %d", len(conf.Qscs))
	for _, qsc := range conf.Qscs {
		log.Tracef("qsc: %s %s", qsc.Name, qsc.NodeAddress)
	}

	log.Info("Starting adapter ports...")
	//启动事件监听 chain node
	w.Add(1)
	go func() {
		startAdapterPorts(conf)
		w.Done()
	}()

	log.Info("Starting qcp consumer...")
	//启动nats 消费
	w.Add(1)
	go func() {
		err = msgqueue.StartQcpConsume(conf)
		if err != nil {
			log.Errorf("Start qcp consume error: %s", err)
			log.Flush()
			os.Exit(1)
		}
		w.Done()
	}()

	go func() {
		w.Wait()
		log.Info("Cassini started ")
	}()
	return
}

func startAdapterPorts(conf *config.Config) {
	// _, err = event.StartEventSubscibe(conf)
	// if err != nil {
	// 	log.Errorf("Start event subscribe error: %s", err)
	// 	log.Flush()
	// 	os.Exit(1)
	// }
	for _, qsc := range conf.Qscs {
		for _, nodeAddr := range strings.Split(qsc.NodeAddress, ",") {
			// go EventsSubscribe(conf.Nats, "tcp://"+nodeAddr, es)
			// subEventFrom += fmt.Sprintf("[%s] ", nodeAddr)
			registerAdapter(nodeAddr, qsc.Name)
		}
	}
}

func registerAdapter(nodeAddr string, chainName string) {
	addrs := strings.Split(nodeAddr, ":")
	if len(addrs) == 2 {
		port, err := strconv.Atoi(addrs[1])
		if err == nil {
			conf := &ports.AdapterConfig{
				ChainName: chainName,
				IP:        addrs[0],
				Port:      port}
			ports.RegisterAdapter(conf)
			return
		}
		log.Errorf("Chain[%s] node address parse error: %s, %v",
			chainName, nodeAddr, err)
	}
	log.Errorf("Adapter ports start error: can not parse chain[%s] node address %s",
		chainName, nodeAddr)
	log.Flush()
	os.Exit(1)

}
