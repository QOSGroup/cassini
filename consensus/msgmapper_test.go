package consensus

import (
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/restclient"
	"github.com/QOSGroup/cassini/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAddMsgToMap(t *testing.T) {

	var events, events1, events2 []types.Event

	events = append(events, newEvent("192.168.1.100:26657", "hashfortest"))
	events = append(events, newEvent("192.168.1.103:26657", "Byzantine"))
	events = append(events, newEvent("192.168.1.101:26657", "hashfortest"))
	events = append(events, newEvent("192.168.1.102:26657", "hashfortest"))

	events1 = append(events1, newEvent("192.168.1.100:26657", "hashfortest"))
	events1 = append(events1, newEvent("192.168.1.103:26657", "hashfortest"))
	events1 = append(events1, newEvent("192.168.1.101:26657", "hashfortest"))
	events1 = append(events1, newEvent("192.168.1.102:26657", "hashfortest"))

	events2 = append(events2, newEvent("192.168.1.100:26657", "hashfortest"))
	events2 = append(events2, newEvent("192.168.1.103:26657", "Byzantine"))
	events2 = append(events2, newEvent("192.168.1.101:26657", "hashfortest"))
	events2 = append(events2, newEvent("192.168.1.102:26657", "Byzantine"))
	events2 = append(events2, newEvent("192.168.1.102:26657", "Byzantine"))

	m := EngineMap{MsgMap: make(map[int64]map[string]string)}
	f := newFerry("fromChain", "toChain", 1)

	for i, event := range events {
		j, err := m.AddMsgToMap(f, event, 3)
		assert.NoError(t, err)
		if i == 3 {
			assert.Equal(t, j, int64(2))
		} else {
			assert.Equal(t, j, int64(0))
		}
	}

	m = EngineMap{MsgMap: make(map[int64]map[string]string)}
	for i, event := range events1 {
		j, err := m.AddMsgToMap(f, event, 3)
		assert.NoError(t, err)
		if i == 2 {
			assert.Equal(t, j, int64(2))
		} else {
			assert.Equal(t, j, int64(0))
		}
	}

	m = EngineMap{MsgMap: make(map[int64]map[string]string)}
	for _, event := range events2 {
		j, err := m.AddMsgToMap(f, event, 3)
		assert.NoError(t, err)
		//fmt.Println(i)
		assert.Equal(t, j, int64(0))
	}
}

func newEvent(node, hash string) types.Event {

	ced := types.CassiniEventDataTx{From: "fromChain", To: "toChain", Height: 1, Sequence: 1}
	ced.HashBytes = []byte(hash)

	event := types.Event{NodeAddress: node, CassiniEventDataTx: ced}
	return event
}

func newFerry(from, to string, sequence int64) *Ferry {
	conf, _ := config.LoadConfig("../config/config.conf")
	f := &Ferry{sequence: 1, conf: conf}
	f.from, f.to = from, to
	f.ConsMap = &ConsensusMap{ConsMap: make(map[int64]map[string]string)}
	f.rmap = make(map[string]*restclient.RestClient)
	for _, node := range strings.Split(conf.GetQscConfig(from).NodeAddress, ",") {
		add := GetAddressFromUrl(node)
		f.rmap[add] = restclient.NewRestClient(node)

	}
	for _, node := range strings.Split(conf.GetQscConfig(to).NodeAddress, ",") {
		add := GetAddressFromUrl(node)
		f.rmap[add] = restclient.NewRestClient(node)

	}

	f.sequence = sequence

	return f
}
