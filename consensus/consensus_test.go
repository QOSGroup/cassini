package consensus

import (
	"testing"

	"github.com/QOSGroup/cassini/config"
	"github.com/stretchr/testify/assert"
)

func TestConsensus32(t *testing.T) {
	conf := &config.Config{}
	conf.ConfigFile = "../cassini.yml"
	err := conf.Load()

	assert.NoError(t, err)

	c := NewConsEngine("qos", "qqs")
	c.F.conf = conf
	N := c.consensus32()
	assert.NotEqual(t, N, 0)

}

func TestGetAddressFromUrl(t *testing.T) {

	assert.Equal(t, GetAddressFromUrl("nats://127.0.0.1"), "127.0.0.1")
	assert.Equal(t, GetAddressFromUrl("tcp://127.0.0.1:26657"), "127.0.0.1:26657")
	assert.Equal(t, GetAddressFromUrl("http://127.0.0.1:8080"), "127.0.0.1:8080")
}
