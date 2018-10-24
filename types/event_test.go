package types

import (
	"github.com/magiconair/properties/assert"
	"github.com/tendermint/tendermint/libs/common"
	"testing"
)

func TestConstructFromTags(t *testing.T) {

	//tags := []common.KVPair{{"qcp.from": "qsc"}, //qsc name 或 qos
	//	{"qcp.to": "qos"}, //qsc name 或 qos
	//	{"qcp.sequence": 18},
	//	{qcp.hash: []byte("hashfortest")}, //TxQcp 做 sha256
	//}

	tags := []common.KVPair{
		{Key: []byte("qcp.from"), Value: []byte("qsc")},
		{Key: []byte("qcp.to"), Value: []byte("qos")},
		{Key: []byte("qcp.sequence"), Value: []byte("10")},
		{Key: []byte("qcp.hash"), Value: []byte("hashfortest")},
	}

	c := CassiniEventDataTx{}
	err := c.ConstructFromTags(tags)
	assert.Equal(t, err, nil)
	assert.Equal(t, c.Sequence, int64(10), "Sequence wrong")
	assert.Equal(t, c.From, "qsc", "From wrong")
	assert.Equal(t, c.To, "qos", "To wrong")
	assert.Equal(t, c.HashBytes, []byte("hashfortest"), "Hashbytes wrong")

}
