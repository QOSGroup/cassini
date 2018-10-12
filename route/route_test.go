package route

import (
	"testing"
	"github.com/QOSGroup/cassini/types"
	)

func TestEvent2queue(t *testing.T) {
	//myroute := route{}
	cEventDatatx := types.CassiniEventDataTx{From: "QSC1",To:"QOS",Sequence:1,HashBytes:[]byte("sha256")}
	event := types.Event{CassiniEventDataTx:cEventDatatx,NodeAddress:}
	Event2queue(&event)
}
