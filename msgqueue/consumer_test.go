package msgqueue

import (
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	"testing"
)

func TestQcpConsume(t *testing.T) {

	//消费消息
	err := make(chan error)
	defer close(err)
	qcpConsume("QSC1", "QOS", DEFAULTSERVERURLS, err)

	assert.Nil(t, err)
}

func TestNATSConsumer_Consume(t *testing.T) {

	i := 0
	cb := func(m *nats.Msg) {
		i++
		tx2 := types.Event{}
		amino.UnmarshalBinary(m.Data, &tx2)
		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, m.Subject, string(m.Data), m.Reply)
		log.Info(tx2.From, tx2.To, tx2.Sequence, string(tx2.HashBytes))
		if string(m.Data) != DEFAULTMSG {
			t.Error("expect the consume msg and the produce msg to match\n")
		}
	}
	//消费消息
	consummer := NATSConsumer{serverUrls: DEFAULTSERVERURLS, subject: DEFAULTSUBJECT, CallBack: cb}
	nc, err := consummer.Connect()
	if err != nil {
		t.Error("couldn't connect to NATS server")
	}
	consummer.Consume(nc)
	//生产消息
	TestNATSProducer_Produce(t)
	select {}
}

//TODO
func TestNATSConsumer_Reply(t *testing.T) {

	i := 0
	cb := func(m *nats.Msg) {
		i++
		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, m.Subject, string(m.Data), m.Reply)
		if string(m.Data) != DEFAULTMSG {
			t.Error("expect the consume msg and the produce msg to match\n")
		}
	}
	//消费消息
	consummer := NATSConsumer{serverUrls: DEFAULTSERVERURLS, subject: DEFAULTSUBJECT, CallBack: cb}
	nc, err := consummer.Connect()
	if err != nil {
		t.Error("couldn't connect to msg server")
	}
	consummer.Reply(nc)
	//生产消息
	TestNATSProducer_ProduceWithReply(t)
	select {}
}

func BenchmarkNATSConsumer_Consume(b *testing.B) {
	i := 0
	cb := func(m *nats.Msg) {
		i += 1
		log.Infof("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
		if string(m.Data) != DEFAULTMSG {
			b.Error("expect the consume msg and the produce msg to match\n")
		}
	}
	consummer := NATSConsumer{serverUrls: DEFAULTSERVERURLS, subject: DEFAULTSUBJECT, CallBack: cb}
	nc, err := consummer.Connect()
	if err != nil {
		b.Error("couldn't connect to msg server")
	}
	consummer.Consume(nc)

	producer := NATSProducer{ServerUrls: DEFAULTSERVERURLS, Subject: DEFAULTSUBJECT}
	np, err := producer.Connect()
	if err != nil {
		b.Error("couldn't connect to msg server")
	}
	for i := 0; i < b.N; i++ { //30000	     51369 ns/op
		producer.Produce(np, []byte(DEFAULTMSG))
	}
}
