package route

import (
	"github.com/huangdao/cassini/types"
		mq "github.com/huangdao/cassini/msgqueue"
	"github.com/tendermint/go-amino"
	"errors"
		"github.com/huangdao/cassini/config"
)

//type route struct{}


func  Event2queue(event *types.Event) error {

	if event == nil || event.HashBytes == nil || event.From == "" || event.To == ""  || event.NodeAddress == ""{
		return errors.New("event is nil")
	}

	eventbytes, _ := amino.MarshalBinary(*event)
	//log.Debug("Event:" , *event)
	//event2 := types.Event{}
	//if amino.UnmarshalBinary(eventbytes,&event2) != nil {
	//	log.Debug("UnmarshalBinary Event:" , event2)
	//}

	subject := event.From + "2" + event.To

	producer := mq.NATSProducer{ServerUrls: config.TestConfig().Nats, Subject: subject}
	np, err := producer.Connect()
	if err != nil {
		return errors.New("couldn't connect to msg server")
	}
	if err := producer.Produce(np, eventbytes); err != nil {
		return err //TODO 错误提示不直接 比如 连接超时
	}

	return nil
}