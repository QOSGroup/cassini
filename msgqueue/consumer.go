//Package msgqueue 从消息队列服务收发消息
package msgqueue

import (
	"errors"
	"fmt"
	"sync"

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

var ce = consensus.NewConsEngine()

var wg sync.WaitGroup

func StartQcpConsume(conf *config.Config) (err error) {

	// qsconfigs := config.DefaultQscConfig()
	qsconfigs := conf.Qscs

	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two qsc names ")
	}

	var subjects string

	es := make(chan error, 1024) //TODO 1024参数按需修改
	defer close(es)

	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			wg.Add(2)
			go qcpConsume(qsconfigs[j].Name, qsconfig.Name, config.DefaultConfig().Nats, es)
			go qcpConsume(qsconfig.Name, qsconfigs[j].Name, config.DefaultConfig().Nats, es)

			subjects += fmt.Sprintf("[%s] [%s]", qsconfigs[j].Name+"2"+qsconfig.Name, qsconfig.Name+"2"+qsconfigs[j].Name)

		}
	}

	wg.Wait()

	if len(es) > 0 {
		for e := range es {
			log.Error(e)
		}
		return errors.New("couldn't start qcp consumer")
	}

	log.Infof("Listening on subjects %s", subjects)

	return
}

//QcpConsumer concume the message from nats server
// from ,to is chain name for example "QOS"
func qcpConsume(from, to, natsServerUrls string, e chan<- error) {

	var i int64 = 0

	defer wg.Add(-1)

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
		e <- err
		return
	}

	if err = consummer.Consume(nc); err != nil {
		e <- err
		return
	}

	return
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

	//nc, err = n.Connect()
	//if err != nil {
	//	return errors.New("the nats.Conn is not available")
	//}

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

	////for test
	//if !strings.Contains(serverUrls, ",") {
	//	log.Debug("serverUrls not contains ','")
	//}
	log.Debugf("connectting to nats :[%s]", serverUrls)

	nc, err = nats.Connect(serverUrls)
	if err != nil {

		log.Errorf("Can't connect %v", err)

		return nil, err
	}
	return
}
