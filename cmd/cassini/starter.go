package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/QOSGroup/cassini/concurrency"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/consensus"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/prometheus"
)

// 命令行 start 命令执行方法
var starter = func() (cancel context.CancelFunc, err error) {

	log.Info("Starting cassini...")

	var w sync.WaitGroup
	errChannel := make(chan error, 1)
	startLog(errChannel)
	startPrometheus(errChannel)
	startEtcd(&w)
	startAdapterPorts(errChannel, &w)
	w.Wait()
	startConsensus(&w)

	go func() {
		w.Wait()
		log.Info("Cassini started ")
	}()
	return
}

func startLog(errChannel <-chan error) {
	go func() {
		for {
			select {
			case e, ok := <-errChannel:
				{
					if ok && e != nil {
						log.Error(e)
					}
				}
			}
		}
	}()
}

func startPrometheus(errChannel chan<- error) {
	log.Info("Starting prometheus exporter...")
	go func() {
		prometheus.StartMetrics(errChannel)
	}()
	log.Info("Prometheus exporter(:39099/metrics) started")
}

func startEtcd(w *sync.WaitGroup) {
	w.Add(1)
	go func() {
		etcd, e := concurrency.StartEmbedEtcd(config.GetConfig())
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
}

func startAdapterPorts(errChannel chan<- error, w *sync.WaitGroup) {
	log.Info("Starting adapter ports...")
	w.Add(1)
	go func() {
		conf := config.GetConfig()
		for _, qsc := range conf.Qscs {
			for _, nodeAddr := range strings.Split(qsc.Nodes, ",") {
				if err := registerAdapter(
					nodeAddr, qsc, errChannel); err != nil {
					errChannel <- err
				}
			}
		}
		w.Done()
	}()
}

func registerAdapter(nodeAddr string, qsc *config.QscConfig,
	errChannel chan<- error) (err error) {
	defer func() {
		if o := recover(); o != nil {
			if err, ok := o.(error); ok {
				errChannel <- fmt.Errorf(
					"Register adapter error: %v", err)
			}
		}
	}()
	addrs := strings.Split(nodeAddr, ":")
	if len(addrs) != 2 {
		err = fmt.Errorf(
			"Adapter ports start error: can not parse chain[%s] node address %s",
			qsc.Name, nodeAddr)
		return
	}
	var port int
	port, err = strconv.Atoi(addrs[1])
	if err != nil {
		err = fmt.Errorf(
			"Chain[%s] node address parse error: %s, %v",
			qsc.Name, nodeAddr, err)
		return
	}
	conf := &ports.AdapterConfig{
		ChainName: qsc.Name,
		ChainType: qsc.Type,
		IP:        addrs[0],
		Port:      port}
	if err = ports.RegisterAdapter(conf); err != nil {
		err = fmt.Errorf(
			"Register adapter error: %v", err)
	}
	return
}

func startConsensus(w *sync.WaitGroup) {
	log.Info("Starting qcp consumer...")
	//启动nats 消费
	w.Add(1)
	go func() {
		err := consensus.StartQcpConsume(config.GetConfig())
		if err != nil {
			log.Errorf("Start qcp consume error: %s", err)
			log.Flush()
			os.Exit(1)
		}
		w.Done()
	}()
}
