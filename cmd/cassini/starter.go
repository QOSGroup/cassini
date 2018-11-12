package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/msgqueue"

	"github.com/etcd-io/etcd/embed"
)

// 命令行 start 命令执行方法
var starter = func(conf *config.Config) (cancel context.CancelFunc, err error) {

	var w sync.WaitGroup
	w.Add(1)
	go func() {
		err = startEtcd(conf.Etcd, &w)
		if err != nil {
			log.Errorf("Start etcd error: %s", err)
		}
	}()
	w.Wait()

	log.Info("Starting cassini...")

	log.Tracef("Qscs: %d", len(conf.Qscs))
	for _, qsc := range conf.Qscs {
		log.Tracef("qsc: %s %s", qsc.Name, qsc.NodeAddress)
	}

	log.Info("Starting event subscriber...")
	//启动事件监听 chain node
	go func() {
		_, err = event.StartEventSubscibe(conf)
		if err != nil {
			log.Errorf("Start event subscribe error: %s", err)
		}
	}()

	log.Info("Starting qcp consumer...")
	//启动nats 消费
	go func() {
		err = msgqueue.StartQcpConsume(conf)
		if err != nil {
			log.Errorf("Start qcp consume error: %s", err)
		}
	}()

	log.Info("Cassini started ")
	return
}

func startEtcd(conf *config.EtcdConfig, w *sync.WaitGroup) (err error) {
	if conf == nil {
		w.Done()
		return
	}
	cfg := embed.NewConfig()

	cfg.ACUrls, cfg.LCUrls, err = ParseUrls(conf.Advertise, conf.Listen)
	cfg.APUrls, cfg.LPUrls, err = ParseUrls(conf.AdvertisePeer, conf.ListenPeer)
	cfg.Dir = fmt.Sprintf("%s.%s", conf.Name, "etcd")
	cfg.InitialCluster = conf.Cluster
	cfg.InitialClusterToken = conf.ClusterToken

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Error("Etcd server start error: ", err)
	}
	w.Done()

	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Info("Etcd server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Info("Etcd server took too long to start!")
	}
	err = <-e.Err()
	log.Error("Etcd running error: ", err)
	return
}

// ParseUrls parse URLs
func ParseUrls(firstURL, secondURL string) (fus, sus []url.URL, err error) {
	var u *url.URL
	u, err = url.Parse(firstURL)
	if err != nil {
		log.Error("Etcd server start error: ", err)
	}
	fus = []url.URL{*u}
	if strings.EqualFold(secondURL, "") {
		sus = []url.URL{*u}
	} else {
		u, err = url.Parse(secondURL)
		if err != nil {
			log.Error("Etcd server start error: ", err)
		}
		sus = []url.URL{*u}
	}
	return
}
