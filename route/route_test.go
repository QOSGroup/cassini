package route

import (
	"testing"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/types"
	"github.com/stretchr/testify/assert"
)

//TODO local nats server
func TestEvent2queue(t *testing.T) {
	conf, err := config.LoadConfig("../config/config.conf")

	cEventDatatx := types.CassiniEventDataTx{From: "QSC1", To: "QOS", Sequence: 1, HashBytes: []byte("sha256")}
	event := types.Event{CassiniEventDataTx: cEventDatatx, NodeAddress: "127.0.0.1:26657"}
	subject, err := Event2queue(conf.Nats, &event)
	assert.Nil(t, err)
	assert.Equal(t, subject, "QSC1"+"2"+"QOS")

	cEventDatatx = types.CassiniEventDataTx{}
	event = types.Event{CassiniEventDataTx: cEventDatatx, NodeAddress: ""}
	subject, err = Event2queue(conf.Nats, &event)
	assert.Equal(t, err.Error(), "event data is empty", "couldn't route empty event")
	assert.Equal(t, subject, "")

	cEventDatatx = types.CassiniEventDataTx{From: "QSC1", To: "", Sequence: 1, HashBytes: []byte("sha256")}
	event = types.Event{CassiniEventDataTx: cEventDatatx, NodeAddress: "127.0.0.1:26657"}
	subject, err = Event2queue(conf.Nats, &event)
	assert.Equal(t, err.Error(), "event data is empty", "couldn't route Incomplete event")
	assert.Equal(t, subject, "")

}
