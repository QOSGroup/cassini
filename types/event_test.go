package types

import (
	"testing"

	"strconv"

	"github.com/QOSGroup/qbase/qcp"
	"github.com/stretchr/testify/assert"
)

func TestConstructFromTags(t *testing.T) {

	tags := make(map[string]string)
	tags[qcp.From] = "qsc"
	tags[qcp.To] = "qos"
	tags[qcp.Sequence] = strconv.FormatInt(7, 10)
	tags[qcp.Hash] = "hashfortest"

	c := CassiniEventDataTx{}
	err := c.ConstructFromTags(tags)
	assert.NoError(t, err)
	assert.Equal(t, c.Sequence, int64(7), "Sequence wrong")
	assert.Equal(t, c.From, "qsc", "From wrong")
	assert.Equal(t, c.To, "qos", "To wrong")
	assert.Equal(t, c.HashBytes, []byte("hashfortest"), "Hashbytes wrong")

}
