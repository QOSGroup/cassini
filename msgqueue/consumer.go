//Package msgqueue 从消息队列服务收发消息
package msgqueue

import (
	"errors"
	// "github.com/QOSGroup/cassini/consensus"
	"github.com/QOSGroup/cassini/log"
	"github.com/nats-io/go-nats"
)

// NATSConsumer Gnatsd consumer
type NATSConsumer struct {
	ServerUrls string //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"

	Subject string //订阅主题

	CallBack func(msg *nats.Msg) //处理消息的回调函数
}

// Connect Connect to gnatsd server
func (n *NATSConsumer) Connect() (nc *nats.Conn, err error) {

	return connect2Nats(n.ServerUrls)
}

// Consume Consume tx msg
func (n *NATSConsumer) Consume(nc *nats.Conn) (err error) {

	if nc == nil {
		return errors.New("the nats.Conn is nil")
	}
	//reconnect to nats server
	i := nc.Status()

	if i != nats.CONNECTED {
		if i != nats.CLOSED {
			nc.Close()
		}
		nc, err = n.Connect()
		if err != nil {
			return errors.New("the nats.Conn is not available")
		}
	}

	//nc, err = n.Connect()
	//if err != nil {
	//	return errors.New("the nats.Conn is not available")
	//}

	subscription, err := nc.Subscribe(n.Subject, n.CallBack)
	if err != nil {
		return errors.New("subscribe failed :" + subscription.Subject)
	}
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Error(err)
	}

	//log.Infof("Listening on [%s]\n", n.subject)

	return nil
}

// Reply Consume tx msg and reply a msg
func (n *NATSConsumer) Reply(nc *nats.Conn) error {

	i := 0
	nc.Subscribe(n.Subject, func(msg *nats.Msg) {
		i++
		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, msg.Subject, string(msg.Data), msg.Reply)
		nc.Publish(msg.Reply, []byte(" thereply"))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Error(err)
	}

	log.Infof("Listening on [%s]\n", n.Subject)

	//runtime.Goexit()
	return nil
}

func connect2Nats(serverUrls string) (nc *nats.Conn, err error) {

	////for test
	//if !strings.Contains(serverUrls, ",") {
	//	log.Debug("serverUrls not contains ','")
	//}
	log.Debugf("connectting to nats [%s]", serverUrls)

	nc, err = nats.Connect(serverUrls)
	if err != nil {

		log.Errorf("Can't connect %v", err)

		return nil, err
	}
	return
}
