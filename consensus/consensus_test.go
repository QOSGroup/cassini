package consensus

import (
	"github.com/QOSGroup/cassini/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsensus32(t *testing.T) {
	conf, _ := config.LoadConfig("../config/config.conf")
	c := NewConsEngine("qos", "qqs")
	c.F.conf = conf
	N := c.consensus32()
	assert.NotEqual(t, N, 0)

}
