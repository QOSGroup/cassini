package concurrency

import (
	"fmt"

	cmn "github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"

	"github.com/etcd-io/etcd/embed"
)

// StartEmbedEtcd start embed etcd
func StartEmbedEtcd(config *config.Config) (etcd *embed.Etcd, err error) {
	if !config.EmbedEtcd || config.Etcd == nil {
		return
	}
	log.Info("Starting etcd...")

	conf := config.Etcd
	cfg := embed.NewConfig()
	cfg.ACUrls, cfg.LCUrls, err = cmn.ParseUrls(conf.Advertise, conf.Listen)
	if err != nil {
		return
	}
	cfg.APUrls, cfg.LPUrls, err = cmn.ParseUrls(conf.AdvertisePeer, conf.ListenPeer)
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
