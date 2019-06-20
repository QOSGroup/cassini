package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParse(t *testing.T) {
	buf := []byte(`{
	"mocks":[
	    {
	        "name": "qos",
	        "rpc": {
		        "nodes": "0.0.0.0:27657"
	        }
	    },
	    {
	        "name": "qsc",
	        "rpc": {
		         "nodes": "0.0.0.0:27658"
	        }
	    }
	]
        	}`)
	conf := GetConfig()
	err := conf.parse(buf)

	assert.NoError(t, err)
	assert.Equal(t, int(2), len(conf.Mocks))
	assert.Equal(t, "qos", conf.Mocks[0].Name)
	assert.Equal(t, "0.0.0.0:27658", conf.Mocks[1].RPC.NodeAddress)
}

func TestLoadConfig(t *testing.T) {
	conf := GetConfig()
	conf.ConfigFile = "./config.conf"
	err := conf.Load()
	assert.NoError(t, err)

	assert.Equal(t, true, conf.Consensus)
	assert.Equal(t, "nats://127.0.0.1:4222", conf.Nats)
	assert.Equal(t, int(2), len(conf.Qscs))
	assert.Equal(t, "fromChain", conf.Qscs[0].Name)
	assert.Equal(t, "toChain", conf.Qscs[1].Name)
	assert.Equal(t, "qos", conf.Qscs[1].Type)
	assert.Equal(t, "127.0.0.1:26657", conf.Qscs[0].NodeAddress)

	assert.Equal(t, int(2), len(conf.Mocks))
	assert.Equal(t, "qos", conf.Mocks[0].Name)
	assert.Equal(t, "0.0.0.0:27657,0.0.0.0:28657", conf.Mocks[1].RPC.NodeAddress)
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
