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
		"listen": "tcp://0.0.0.0:27657"
	        }
	    },
	    {
	        "name": "qsc",
	        "rpc": {
		"listen": "tcp://0.0.0.0:27658"
	        }
	    }
	]
        	}`)
	conf, err := CreateConfig(buf)
	assert.NoError(t, err)
	assert.Equal(t, int(2), len(conf.Mocks))
	assert.Equal(t, "qos", conf.Mocks[0].Name)
	assert.Equal(t, "tcp://0.0.0.0:27658", conf.Mocks[1].RPC.ListenAddress)
}

func TestTransfrom(t *testing.T) {
	var a uint32
	a--
	t.Logf("a: %v", a)
	assert.Equal(t, true, true)
}
