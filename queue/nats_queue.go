package queue

import (
	"errors"
	"fmt"

	"github.com/QOSGroup/cassini/log"
	"github.com/nats-io/go-nats"
)

// NatsQueue wraps nats client as a message queue service
type NatsQueue struct {
	Subject string
	Config  string
}

// Init message queue
func (q *NatsQueue) Init() error {
	return nil
}

// NewProducer returns a new producer for the message queue
func (q *NatsQueue) NewProducer() (p Producer, err error) {
	conn, err := connect2Nats(q.Config)
	if err != nil {
		log.Errorf("Connect error: %v", err)
		return nil, err
	}
	return &NatsProducer{queue: q, conn: conn}, nil
}

// NewComsumer returns a new comsumer for the message queue
func (q *NatsQueue) NewComsumer() (c Comsumer, err error) {
	conn, err := connect2Nats(q.Config)
	if err != nil {
		log.Errorf("Connect error: %v", err)
		return nil, err
	}
	return &NatsComsumer{queue: q, conn: conn}, nil
}

// NatsProducer define the producer for local message queue based on channel
type NatsProducer struct {
	queue *NatsQueue
	conn  *nats.Conn
}

// Produce sends data to local message queue
func (p *NatsProducer) Produce(data []byte) (err error) {
	if p.conn == nil {
		return fmt.Errorf("the nats.Conn is nil - %s, %s", p.Config(), p.Subject())
	}

	//reconnect to nats server
	i := p.conn.Status()
	if i != nats.CONNECTED {

		if i != nats.CLOSED {
			p.conn.Close()
		} //status==2 closed

		p.conn, err = connect2Nats(p.Config())
		if err != nil {
			log.Errorf("Reconnect error: %v", err)
			return errors.New("the nats.Conn is not available")
		}
	}

	if e := p.conn.Publish(p.Subject(), data); e != nil {
		log.Errorf("Send data error: %v", err)
		return errors.New("send event to nats server faild")
	}

	p.conn.Flush()

	if err := p.conn.LastError(); err != nil {
		log.Error(err)
	}

	return nil
}

// Subject returns subject of local message queue
func (p *NatsProducer) Subject() string {
	return p.queue.Subject
}

// Config returns config of local message queue
func (p *NatsProducer) Config() string {
	return p.queue.Config
}

// NatsComsumer define the comsumer for local message queue based on channel
type NatsComsumer struct {
	queue *NatsQueue
	conn  *nats.Conn
}

// Subscribe sets the listener for local message queue based on channel
func (c *NatsComsumer) Subscribe(listener Listener) (err error) {
	if c.conn == nil {
		return fmt.Errorf("the nats.Conn is nil - %s, %s", c.Config(), c.Subject())
	}
	//reconnect to nats server
	i := c.conn.Status()

	if i != nats.CONNECTED {
		if i != nats.CLOSED {
			c.conn.Close()
		}
		c.conn, err = connect2Nats(c.Config())
		if err != nil {
			log.Errorf("Reconnect error: %v", err)
			return errors.New("the nats.Conn is not available")
		}
	}

	msgHandler := func(msg *nats.Msg) {
		listener(msg.Data, c)
	}

	subscription, err := c.conn.Subscribe(c.Subject(), msgHandler)
	if err != nil {
		return errors.New("subscribe failed :" + subscription.Subject)
	}
	c.conn.Flush()

	if err := c.conn.LastError(); err != nil {
		log.Error(err)
	}
	return nil
}

// Subject returns subject of local message queue
func (c *NatsComsumer) Subject() string {
	return c.queue.Subject
}

// Config returns config of local message queue
func (c *NatsComsumer) Config() string {
	return c.queue.Config
}

func connect2Nats(conf string) (nc *nats.Conn, err error) {
	// log.Debugf("Connect to nats [%s]", conf)

	nc, err = nats.Connect(conf)
	if err != nil {
		return nil, err
	}
	return
}
