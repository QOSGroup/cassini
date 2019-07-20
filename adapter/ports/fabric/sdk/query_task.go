package sdk

import (
	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk/utils"
	"github.com/QOSGroup/cassini/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// QueryTask is the query task
type QueryTask struct {
	channelClient *channel.Client
	targets       []fab.Peer
	id            string
	ccID          string
	args          *Args
	callback      func(err error)
	verbose       bool
	payloadOnly   bool
}

// NewQueryTask creates a new query Task
func NewQueryTask(id string, channelClient *channel.Client, targets []fab.Peer,
	chaincodeID string, args *Args,
	verbose bool, payloadOnly bool, callback func(err error)) *QueryTask {
	return &QueryTask{
		id:            id,
		channelClient: channelClient,
		targets:       targets,
		ccID:          chaincodeID,
		args:          args,
		callback:      callback,
		verbose:       verbose,
		payloadOnly:   payloadOnly,
	}
}

// Invoke invokes the query task
func (t *QueryTask) Invoke() {
	var opts []channel.RequestOption
	if len(t.targets) > 0 {
		opts = append(opts, channel.WithTargets(t.targets...))
	}
	if response, err := t.channelClient.Query(
		channel.Request{
			ChaincodeID: t.ccID,
			Fcn:         t.args.Func,
			Args:        utils.AsBytes(t.args.Args),
		},
		opts...,
	); err != nil {
		log.Debugf("(%s) - Error querying chaincode: %s\n", t.id, err)
		t.callback(err)
	} else {
		log.Debugf("(%s) - Chaincode query was successful\n", t.id)

		log.Debugf("response payload: %s", response.Payload)

		t.callback(nil)
	}
}
