package route

import (
	"errors"

	"github.com/QOSGroup/cassini/log"
	mq "github.com/QOSGroup/cassini/msgqueue"
	"github.com/QOSGroup/cassini/types"
	"github.com/tendermint/go-amino"
)

//type route struct{}

// Event2queue produce event to message queue (Nats)
func Event2queue(nats string, event *types.Event) (subject string, err error) {

	if event == nil {

		return "", errors.New("event is nil")
	}
	// log.Infof("event from: %s, to: %s, nodes: %s, sequence: %d, hash: %v",
	// 	event.From, event.To, event.NodeAddress, event.Sequence, event.HashBytes)
	if event.HashBytes == nil || event.From == "" || event.To == "" || event.NodeAddress == "" {
		return "", errors.New("event data is empty")
	}

	eventbytes, _ := amino.MarshalBinaryLengthPrefixed(*event)

	subject = event.From + "2" + event.To

	producer := mq.NATSProducer{ServerUrls: nats, Subject: subject}

	np, err := producer.Connect() //TODO don't connect every time

	if err != nil {

		return "", errors.New("couldn't connect to msg server")
	}

	defer np.Close()

	if err := producer.Produce(np, eventbytes); err != nil {
		return "", err
	}

	log.Infof("routed event from[%s] sequence[#%d] to subject [%s] ", event.NodeAddress, event.Sequence, subject)

	return subject, nil
}
