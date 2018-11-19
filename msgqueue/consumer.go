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
	"sync"
)

var wg sync.WaitGroup

// StartQcpConsume Start to consume tx msg
func StartQcpConsume(conf *config.Config) (err error) {

	// qsconfigs := config.DefaultQscConfig()
	qsconfigs := conf.Qscs

	if len(qsconfigs) < 2 {
		return errors.New("config error , at least two chain names ")
	}

	var subjects string

	es := make(chan error, 1024) //TODO 1024参数按需修改
	defer close(es)

	engines := make([]*consensus.ConsEngine, 0)
	var ce *consensus.ConsEngine
	for i, qsconfig := range qsconfigs {
		for j := i + 1; j < len(qsconfigs); j++ {
			wg.Add(2)

			ce = createConsEngine(qsconfigs[j].Name, qsconfig.Name, conf, es)
			engines = append(engines, ce)
			ce = createConsEngine(qsconfig.Name, qsconfigs[j].Name, conf, es)
			engines = append(engines, ce)

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

	log.Infof("listening on subjects %s", subjects)

	for _, ce := range engines {
		go ce.StartEngine()
		go ce.F.StartFerry()
	}

	return
}

func createConsEngine(from, to string, conf *config.Config, e chan<- error) (ce *consensus.ConsEngine) {
	ce = consensus.NewConsEngine(from, to)

	seq, err := ce.F.GetSequenceFromChain(from, to, "in") // seq= toChain's in/fromchain/maxseq
	if err != nil {
		log.Errorf("Create consensus engine error: %v", err)
	} else {
		log.Debugf("Create consensus engine query chain %s in-sequence: %d", to, seq)
		ce.SetSequence(from, to, seq)
	}

	go qcpConsume(ce, from, to, conf, e)
	return ce
}

//QcpConsumer consume the message from nats server
//
// from ,to is chain name for example "QOS"
func qcpConsume(ce *consensus.ConsEngine, from, to string, conf *config.Config, e chan<- error) {
	log.Debugf("Consume qcp from [%s] to [%s]", from, to)

	var i int64

	defer wg.Done()

	cb := func(m *nats.Msg) {

		i++

		tx := types.Event{}
		amino.UnmarshalBinary(m.Data, &tx)

		log.Infof("[#%d] Consume subject [%s] sequence [#%d] nodeAddress '%s'", i, m.Subject, tx.Sequence, tx.NodeAddress)

		// 监听到交易事件后立即查询需要等待一段时间才能查询到交易数据；
		//TODO 优化
		// 需要监听下一个块的New Block 事件以确认交易数据入块，abci query 接口才能够查询出交易；
		// 同时提供定时出发机制，以保证共识模块在交易事件丢失或网络错误等问题出现时仍然能够正常运行。
		//if conf.EventWaitMillitime > 0 {
		//	time.Sleep(time.Duration(conf.EventWaitMillitime) * time.Millisecond)
		//}

		ce.Add2Engine(m)
	}

	consummer := NATSConsumer{serverUrls: conf.Nats, subject: from + "2" + to, CallBack: cb}

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

// NATSConsumer Gnatsd consumer
type NATSConsumer struct {
	serverUrls string //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"

	subject string //订阅主题

	CallBack func(msg *nats.Msg) //处理消息的回调函数
}

// Connect Connect to gnatsd server
func (n *NATSConsumer) Connect() (nc *nats.Conn, err error) {

	return connect2Nats(n.serverUrls)
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

// Reply Consume tx msg and reply a msg
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
	log.Debugf("connectting to nats [%s]", serverUrls)

	nc, err = nats.Connect(serverUrls)
	if err != nil {

		log.Errorf("Can't connect %v", err)

		return nil, err
	}
	return
}
