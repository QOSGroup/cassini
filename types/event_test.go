package types

import (
	"testing"

	"github.com/QOSGroup/qbase/qcp"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
	"strconv"
)

func TestConstructFromTags(t *testing.T) {

	tags := []common.KVPair{
		{Key: []byte(qcp.QcpFrom), Value: []byte("qsc")},
		{Key: []byte(qcp.QcpTo), Value: []byte("qos")},
		{Key: []byte(qcp.QcpSequence), Value: []byte(strconv.FormatInt(7, 10))},
		{Key: []byte(qcp.QcpHash), Value: []byte("hashfortest")},
	}

	c := CassiniEventDataTx{}
	err := c.ConstructFromTags(tags)
	assert.NoError(t, err)
	assert.Equal(t, c.Sequence, int64(7), "Sequence wrong")
	assert.Equal(t, c.From, "qsc", "From wrong")
	assert.Equal(t, c.To, "qos", "To wrong")
	assert.Equal(t, c.HashBytes, []byte("hashfortest"), "Hashbytes wrong")

}
