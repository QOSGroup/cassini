package msgqueue

import (
	"testing"

	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/consensus"
	"github.com/QOSGroup/cassini/log"
	"github.com/nats-io/go-nats"
)

//func TestQcpConsume(t *testing.T) {
//
//	conf, _ := config.LoadConfig("../config/config.conf")
//	//消费消息
//	err := make(chan error)
//	defer close(err)
//	ce := newConsEngine("qqs", "qos")
//	wg.Add(1)
//	qcpConsume(ce, "qqs", "qos", conf, err)
//	wg.Wait()
//	//StartQcpConsume(config.TestConfig())
//
//	assert.Nil(t, err)
//}

//func TestNATSConsumer_Consume(t *testing.T) {
//
//	i := 0
//	cb := func(m *nats.Msg) {
//		i++
//		tx2 := types.Event{}
//		amino.UnmarshalBinary(m.Data, &tx2)
//		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, m.Subject, string(m.Data), m.Reply)
//		log.Info(tx2.From, tx2.To, tx2.Sequence, string(tx2.HashBytes))
//		if string(m.Data) != DEFAULTMSG {
//			t.Error("expect the consume msg and the produce msg to match\n")
//		}
//	}
//
//	//go TestNATSProducer_Produce(t)
//	//消费消息
//	consummer := NATSConsumer{serverUrls: DEFAULTNATSURLS, subject: DEFAULTSUBJECT, CallBack: cb}
//	nc, err := consummer.Connect()
//	if err != nil {
//		t.Error("couldn't connect to NATS server")
//	}
//	consummer.Consume(nc)
//
//}

//func TestNATSConsumer_Reply(t *testing.T) {
//
//	i := 0
//	cb := func(m *nats.Msg) {
//		i++
//		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, m.Subject, string(m.Data), m.Reply)
//		if string(m.Data) != DEFAULTMSG {
//			t.Error("expect the consume msg and the produce msg to match\n")
//		}
//	}
//	//消费消息
//	consummer := NATSConsumer{serverUrls: DEFAULTNATSURLS, subject: DEFAULTSUBJECT, CallBack: cb}
//	nc, err := consummer.Connect()
//	if err != nil {
//		t.Error("couldn't connect to msg server")
//	}
//	consummer.Reply(nc)
//	//生产消息
//	TestNATSProducer_ProduceWithReply(t)
//	select {}
//}

func BenchmarkNATSConsumer_Consume(b *testing.B) {
	i := 0
	cb := func(m *nats.Msg) {
		i++
		log.Infof("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
		if string(m.Data) != DEFAULTMSG {
			b.Error("expect the consume msg and the produce msg to match\n")
		}
	}
	consummer := NATSConsumer{serverUrls: DEFAULTNATSURLS, subject: DEFAULTSUBJECT, CallBack: cb}
	nc, err := consummer.Connect()
	if err != nil {
		b.Error("couldn't connect to msg server")
	}
	consummer.Consume(nc)

	producer := NATSProducer{ServerUrls: DEFAULTNATSURLS, Subject: DEFAULTSUBJECT}
	np, err := producer.Connect()
	if err != nil {
		b.Error("couldn't connect to msg server")
	}
	for i := 0; i < b.N; i++ { //30000	     51369 ns/op
		producer.Produce(np, []byte(DEFAULTMSG))
	}
}

func newConsEngine(from, to string) *consensus.ConsEngine {

	conf, _ := config.LoadConfig("../config/config.conf")
	ce := new(consensus.ConsEngine)
	ce.M = &consensus.EngineMap{MsgMap: make(map[int64]map[string]string)}
	ce.F = consensus.NewFerry(conf, from, to, 0)

	ce.Setfrom(from)
	ce.Setto(to)
	return ce
}
