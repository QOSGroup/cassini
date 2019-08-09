package config

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	viper.SetConfigFile("./../cassini.yml")
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	conf := GetConfig()
	err = conf.Load()
	assert.NoError(t, err)

	assert.Equal(t, true, conf.Consensus)
	assert.Equal(t, "nats://127.0.0.1:4222", conf.Queue)
	assert.Equal(t, int64(2000), conf.EventWaitMillitime)
	assert.Equal(t, true, conf.UseEtcd)
	assert.Equal(t, "etcd://127.0.0.1:2379", conf.Lock)
	assert.Equal(t, int64(5000), conf.LockTTL)
	assert.Equal(t, true, conf.EmbedEtcd)

	assert.Equal(t, int(2), len(conf.Qscs))
	assert.Equal(t, "fromChain", conf.Qscs[0].Name)
	assert.Equal(t, "qstars", conf.Qscs[0].Type)
	assert.Equal(t, "127.0.0.1:26657", conf.Qscs[0].Nodes)
	assert.Equal(t, "toChain", conf.Qscs[1].Name)
	assert.Equal(t, "qos", conf.Qscs[1].Type)
	assert.Equal(t, "127.0.0.1:27657", conf.Qscs[1].Nodes)

	assert.Equal(t, "dev-cassini", conf.Etcd.Name)
	assert.Equal(t, "http://127.0.0.1:2379", conf.Etcd.Advertise)
	assert.Equal(t, "http://127.0.0.1:2380", conf.Etcd.AdvertisePeer)
	assert.Equal(t, "dev-cassini-cluster", conf.Etcd.ClusterToken)
	assert.Equal(t, "dev-cassini=http://127.0.0.1:2380", conf.Etcd.Cluster)
}

func TestTransfrom(t *testing.T) {
	var a uint32
	a--
	t.Logf("a: %v", a)
	assert.Equal(t, true, true)
}

//func TestNon(t *testing.T) {
//	conf, err := LoadConfig("./dev_config.conf")
//	fmt.Println(err)
//	s := conf.Test
//	//var s string
//	//s = "c1\000"
//	//s += "c2\000"
//	//s += "c3\000"
//	fmt.Println(s)
//	ss := strings.Split(s, "\000")
//	fmt.Println(len(ss))
//	assert.Equal(t, 3, strings.Count(s, "\000"))
//
//}
