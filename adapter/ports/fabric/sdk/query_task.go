package sdk

import (
	"github.com/QOSGroup/cassini/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/action"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/chaincode/utils"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/printer"
)

// QueryTask is the query task
type QueryTask struct {
	channelClient *channel.Client
	targets       []fab.Peer
	id            string
	ccID          string
	args          *action.ArgStruct
	callback      func(err error)
	printer       printer.Printer
	verbose       bool
	payloadOnly   bool
}

// NewQueryTask creates a new query Task
func NewQueryTask(id string, channelClient *channel.Client, targets []fab.Peer,
	chaincodeID string, args *action.ArgStruct,
	printer printer.Printer, verbose bool, payloadOnly bool, callback func(err error)) *QueryTask {
	return &QueryTask{
		id:            id,
		channelClient: channelClient,
		targets:       targets,
		ccID:          chaincodeID,
		args:          args,
		callback:      callback,
		printer:       printer,
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

		if t.verbose {
			t.printer.PrintTxProposalResponses(response.Responses, t.payloadOnly)
		}

		t.callback(nil)
	}
}
