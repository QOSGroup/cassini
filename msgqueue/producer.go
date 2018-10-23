package msgqueue

import (
	"errors"
	"github.com/QOSGroup/cassini/log"
	"github.com/nats-io/go-nats"
	"time"
)

//type Producer interface {
//	Produce(nc *nats.Conn, msg []byte) error
//}

type NATSProducer struct {
	ServerUrls string //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"

	Subject string //主题
}

func (n *NATSProducer) Connect() (nc *nats.Conn, err error) {

	return connect2Nats(n.ServerUrls)
}

func (n *NATSProducer) Produce(nc *nats.Conn, msg []byte) (err error) {

	//if nc == nil {
	//	return errors.New("the nats.Conn is nil")
	//}
	//
	////reconnect to nats server
	//i := nc.Status()
	//if i != nats.CONNECTED {
	//
	//	if i != nats.CLOSED {
	//		nc.Close()
	//	} //status==2 closed
	//
	//	nc, err = n.Connect()
	//	if err != nil {
	//
	//		return errors.New("the nats.Conn is not available")
	//	}
	//}
	nc, err = n.Connect()
	if err != nil {
		return errors.New("the nats.Conn is not available")
	}
	if e := nc.Publish(n.Subject, msg); e != nil {

		return errors.New("send event to nats server faild")
	}

	nc.Flush()

	if err := nc.LastError(); err != nil {

		log.Error(err)
	}

	return nil
}

//TODO
func (n *NATSProducer) ProduceWithReply(nc *nats.Conn, reply string, payload []byte) error {

	msg, err := nc.Request(n.Subject, payload, 100*time.Millisecond)
	if err != nil {
		if nc.LastError() != nil {
			log.Errorf("Error in Request: %v", nc.LastError())
		}
		log.Errorf("Error in Request: %v", err)
	}

	log.Infof("Published [%s] : '%s'", n.Subject, payload)
	log.Infof("Received [%v] : '%s'", msg.Subject, string(msg.Data))
	log.Infof("Reply [%v] : '%s'", msg.Subject, string(msg.Reply))
	return nil
}
