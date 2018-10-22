//Package msgqueue 从消息队列服务收发消息
package msgqueue

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/consensus"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/nats-io/go-nats"
	"github.com/tendermint/go-amino"
)

//type Consumer interface{
//	Consume(nc *nats.Conn) (err error)
//}

type QcpConsumer struct {
	NATSConsumer
}

var ce = consensus.NewConsEhgine()

func StartQcpConsume(conf *config.Config) (err error) {

	qsconfigs := config.DefaultQscConfig()

	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two qsc names ")
	}

	var subjects string
	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			go QcpConsume(qsconfigs[j].Name, qsconfig.Name, config.DefaultConfig().Nats)
			go QcpConsume(qsconfig.Name, qsconfigs[j].Name, config.DefaultConfig().Nats)

			subjects += fmt.Sprintf("[%s] [%s]", qsconfigs[j].Name+"2"+qsconfig.Name, qsconfig.Name+"2"+qsconfigs[j].Name)
			//err = QcpConsume("QSC1", "QOS", config.DefaultConfig().Nats) //TODO

			//if err == nil {
			//	log.Infof("Listening on subject [%s]", "QSC12QOS")
			//}
		}
	}

	log.Infof("Listening on subjects %s", subjects)

	return
}

//QcpConsumer concume the message from nats server
// from ,to is chain name for example "QOS"
func QcpConsume(from, to, natsServerUrls string) error {

	var i int64 = 0

	cb := func(m *nats.Msg) {

		i++

		tx := types.Event{}
		amino.UnmarshalBinary(m.Data, &tx)

		log.Infof("[#%d] Consume subject [%s] sequence [#%d] nodeAddress '%s'", i, m.Subject, tx.Sequence, tx.NodeAddress)

		ce.Add2Engine(m)
	}

	consummer := NATSConsumer{serverUrls: natsServerUrls, subject: from + "2" + to, CallBack: cb}

	nc, err := consummer.Connect()
	if err != nil {
		return errors.New("couldn't connect to NATS server")
	}

	if err = consummer.Consume(nc); err != nil {
		return err
	}

	return nil
}

type NATSConsumer struct {
	serverUrls string //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"

	subject string //订阅主题

	CallBack func(msg *nats.Msg) //处理消息的回调函数
}

func (n *NATSConsumer) Connect() (nc *nats.Conn, err error) {

	return connect2Nats(n.serverUrls)
}

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

	subscription, err := nc.Subscribe(n.subject, n.CallBack)
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

func (n *NATSConsumer) Reply(nc *nats.Conn) error {

	i := 0
	nc.Subscribe(n.subject, func(msg *nats.Msg) {
		i++
		log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, msg.Subject, string(msg.Data), msg.Reply)
		nc.Publish(msg.Reply, []byte(" thereply"))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Error(err)
	}

	log.Infof("Listening on [%s]\n", n.subject)

	//runtime.Goexit()
	return nil
}

func connect2Nats(serverUrls string) (nc *nats.Conn, err error) {

	nc, err = nats.Connect(serverUrls)
	//log.Debugf("connectting to nats :[%s]", serverUrls)
	if err != nil {

		log.Errorf("Can't connect: %v", err)

		return nil, err
	}
	return
}
