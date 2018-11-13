package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
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

	log.Info("Starting cassini...")

	var w sync.WaitGroup
	w.Add(1)
	go func() {
		var etcd *embed.Etcd
		etcd, err = startEtcd(conf)
		if err != nil {
			log.Error("Etcd server start error: ", err)
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
		err = <-etcd.Err()
		log.Error("Etcd running error: ", err)
	}()

	log.Tracef("Qscs: %d", len(conf.Qscs))
	for _, qsc := range conf.Qscs {
		log.Tracef("qsc: %s %s", qsc.Name, qsc.NodeAddress)
	}

	log.Info("Starting event subscriber...")
	//启动事件监听 chain node
	w.Add(1)
	go func() {
		_, err = event.StartEventSubscibe(conf)
		if err != nil {
			log.Errorf("Start event subscribe error: %s", err)
			os.Exit(1)
		}
		w.Done()
	}()

	log.Info("Starting qcp consumer...")
	//启动nats 消费
	w.Add(1)
	go func() {
		err = msgqueue.StartQcpConsume(conf)
		if err != nil {
			log.Errorf("Start qcp consume error: %s", err)
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

func startEtcd(config *config.Config) (etcd *embed.Etcd, err error) {
	if config.UseEtcd || config.Etcd == nil {
		return
	}
	log.Info("Starting etcd...")

	conf := config.Etcd
	cfg := embed.NewConfig()
	cfg.ACUrls, cfg.LCUrls, err = ParseUrls(conf.Advertise, conf.Listen)
	if err != nil {
		return
	}
	cfg.APUrls, cfg.LPUrls, err = ParseUrls(conf.AdvertisePeer, conf.ListenPeer)
	if err != nil {
		return
	}
	cfg.Dir = fmt.Sprintf("%s.%s", conf.Name, "etcd")
	cfg.InitialCluster = conf.Cluster
	cfg.InitialClusterToken = conf.ClusterToken
	cfg.Name = conf.Name
	cfg.Debug = false

	return embed.StartEtcd(cfg)
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
