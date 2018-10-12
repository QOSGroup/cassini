//Package msgqueue 从消息队列服务收发消息
package msgqueue

import (
	"github.com/nats-io/go-nats"
	"github.com/huangdao/cassini/log"
	"errors"
	"github.com/huangdao/cassini/types"
	"github.com/tendermint/go-amino"
	"github.com/huangdao/cassini/consensus"
	"github.com/huangdao/cassini/config"
	)

//type Consumer interface{
//	Consume(nc *nats.Conn) (err error)
//}

type QcpConsumer struct{

	NATSConsumer

}

func StartQcpConsume(conf *config.Config) ( err error) {

	qsconfigs := config.TestQscConfig()
	if len(qsconfigs) < 2 {
		return errors.New("config error , at leat two qsc names ")
	}

	for i,qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs) ; j++ {
			err =QcpConsume(qsconfig.Name,qsconfigs[j].Name,config.DefaultConfig().Nats)
			err =QcpConsume(qsconfigs[j].Name,qsconfig.Name,config.DefaultConfig().Nats)
		}
	}

	return
}

//QcpConsumer concume the message from nats server
// from ,to is chain name for example "QOS"
func  QcpConsume(from,to,natsServerUrls string)  error{

	cb := qcpCallBack

	consummer := NATSConsumer{serverUrls: natsServerUrls, subject: from + "2" +to, CallBack: cb}

	nc, err := consummer.Connect()
	if err != nil {
		return errors.New("couldn't connect to NATS server")
	}

	if err = consummer.Consume(nc) ; err != nil {
		return err
	}

	return nil
}

func qcpCallBack(m *nats.Msg) {
	i := 0
	tx2 := types.Event{}
	amino.UnmarshalBinary(m.Data,&tx2)
	log.Infof("[#%d] Received on [%s]: '%s' Relpy:'%s'\n", i, m.Subject,string(m.Data), m.Reply)
	log.Info(tx2.From ,tx2.To,tx2.Sequence,string(tx2.HashBytes))

	msgMapper := consensus.MsgMapper{}
	msgMapper.AddMsgToMap(m)

}

type NATSConsumer struct {
	serverUrls string              //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"
	subject    string              //订阅主题
	CallBack   func(msg *nats.Msg) //处理消息的回调函数
}

func (n *NATSConsumer) Connect() (nc *nats.Conn, err error) {
	nc, err = nats.Connect(n.serverUrls)
	if err != nil {
		log.Error("Can't connect: %v\n", err)
		return nil, err
	}
	return
}

func (n *NATSConsumer) Consume(nc *nats.Conn) (err error) {

	if nc == nil {
		return errors.New("the nats.Conn is nil")
	}
	//reconnect to nats server
	i := nc.Status()

	if i != nats.CONNECTED {
		if i != nats.CLOSED {nc.Close()}
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

	log.Infof("Listening on [%s]\n", n.subject)

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