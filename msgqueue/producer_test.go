package msgqueue

import (
	"testing"
)

const (
	DEFAULTMSG      = "the test msg"
	DEFAULTNATSURLS = "nats://127.0.0.1:4222"
	DEFAULTSUBJECT  = "QSC12QOS"
)

func TestNATSProducer_Produce(t *testing.T) {
	producer := NATSProducer{ServerUrls: DEFAULTNATSURLS, Subject: DEFAULTSUBJECT}
	np, err := producer.Connect()
	if err != nil {
		t.Error("couldn't connect to msg server")
	}
	if err := producer.Produce(np, []byte(DEFAULTMSG)); err != nil {
		t.Error(err)
	}
}

//func TestNATSProducer_ProduceWithReply(t *testing.T) {
//	producer := NATSProducer{ServerUrls: DEFAULTNATSURLS, Subject: DEFAULTSUBJECT}
//	np, err := producer.Connect()
//	if err != nil {
//		t.Error("couldn't connect to msg server")
//	}
//	if err := producer.ProduceWithReply(np, "reply test", []byte(DEFAULTMSG)); err != nil {
//		t.Error(err)
//	}
//}

func BenchmarkNATSProducer_Produce(b *testing.B) {
	producer := NATSProducer{ServerUrls: DEFAULTNATSURLS, Subject: DEFAULTSUBJECT}
	np, err := producer.Connect()
	if err != nil {
		b.Error("couldn't connect to msg server")
	}
	for i := 0; i < b.N; i++ { //1000	   1.841105 ms/op
		producer.Produce(np, []byte(DEFAULTMSG))
	}
}
