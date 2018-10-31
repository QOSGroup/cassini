package event

//
//import (
//	"context"
//	"fmt"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/QOSGroup/cassini/config"
//	"github.com/QOSGroup/cassini/mock"
//	"github.com/stretchr/testify/assert"
//	tmtypes "github.com/tendermint/tendermint/types"
//)
//
//func TestEventssubRemote(t *testing.T) {
//
//	// mockergit
//	mc := config.TestQscMockConfig()
//	cancelMock, err := mock.StartMock(*mc)
//	defer cancelMock()
//
//	addr := mc.RPC.NodeAddress
//
//	ipPort := strings.SplitN(addr, ":", 2)
//	if len(ipPort) != 2 {
//		err = fmt.Errorf("Ip and port parse error: %v", addr)
//		assert.NoError(t, err)
//	}
//	//go EventSubscribe("tcp://192.168.168.27:26657")
//	// EventSubscribe(fmt.Sprintf("tcp://127.0.0.1:%v", ipPosrt[1]))
//	remote := fmt.Sprintf("tcp://127.0.0.1:%v", ipPort[1])
//	t.Log("remote", remote)
//	txs := make(chan interface{})
//	var cancel context.CancelFunc
//	cancel, err = SubscribeRemote(remote, "test-client", "tm.event = 'Tx'", txs)
//
//	assert.NoError(t, err)
//	// t.Error(err)
//	defer cancel()
//
//	done := make(chan struct{})
//	go func() {
//		i := 0
//		for e := range txs {
//			i++
//			eventData := e.(tmtypes.EventDataTx)
//			t.Log("Tx Height: ", eventData.Height, ", ", i)
//
//			//cassiniEventDataTx := cassinitypes.CassiniEventDataTx{}
//			//
//			//for _, tag := range eventData.Result.Tags {
//			//	if string(tag.Key) == "qcp.from" {
//			//		cassiniEventDataTx.From = string(tag.Value)
//			//	}
//			//	if string(tag.Key) == "qcp.to" {
//			//		cassiniEventDataTx.To = string(tag.Value)
//			//	}
//			//	if string(tag.Key) == "qcp.hash" {
//			//		cassiniEventDataTx.HashBytes = tag.Value
//			//	}
//			//	if string(tag.Key) == "qcp.sequence" {
//			//		cassiniEventDataTx.Sequence , err =  strconv.ParseInt(string(tag.Value), 10, 64)
//			//		if err != nil { t.Error("sequence not correct")}
//			//		//bin_buf := bytes.NewBuffer(tag.Value)
//			//		//binary.Read(bin_buf, binary.BigEndian, &cassiniEventDataTx.Sequence)
//			//	}
//			//}
//			//fmt.Println(cassiniEventDataTx)
//
//			if i >= 10 {
//				close(done)
//			}
//		}
//		assert.Equal(t, int(10), i)
//	}()
//
//	select {
//	case <-done:
//	case <-time.After(10 * time.Second):
//		t.Fatal("did not receive a transaction after 10 sec.")
//	}
//}
