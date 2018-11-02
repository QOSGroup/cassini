package tx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTxMock(t *testing.T) {

	tx := NewTxQcpMock("from", "to", 99, 199)

	assert.Equal(t, tx.From, "from")
	assert.Equal(t, tx.To, "to")
	assert.Equal(t, tx.BlockHeight, int64(99))
	assert.Equal(t, tx.Sequence, int64(199))

}
