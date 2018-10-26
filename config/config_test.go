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
	conf, err := CreateConfig(buf)
	assert.NoError(t, err)
	assert.Equal(t, int(2), len(conf.Mocks))
	assert.Equal(t, "qos", conf.Mocks[0].Name)
	assert.Equal(t, "0.0.0.0:27658", conf.Mocks[1].RPC.NodeAddress)
}

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig("./config.conf")
	assert.NoError(t, err)

	assert.Equal(t, "no", conf.Consensus)
	assert.Equal(t, "nats://192.168.1.99:4222", conf.Nats)
	assert.Equal(t, int(2), len(conf.Qscs))
	assert.Equal(t, "qqs", conf.Qscs[1].Name)
	assert.Equal(t, "192.168.1.100:26657", conf.Qscs[0].NodeAddress)

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
