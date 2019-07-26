package sdk

import (
	"github.com/QOSGroup/cassini/adapter/ports/fabric/sdk/utils"
	"github.com/QOSGroup/cassini/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

type queryAction struct {
	Action
	numInvoked uint32
	done       chan bool
}

func newQueryAction() (*queryAction, error) {
	action := &queryAction{done: make(chan bool)}
	err := action.Initialize()
	return action, err
}

func (a *queryAction) query(channelID, chaincodeID string, argsArray []Args) (string, error) {
	log.Debug("queryAction.query")
	a.Set(channelID, chaincodeID, argsArray)
	channelClient, err := a.ChannelClient()
	if err != nil {
		return "", errors.Errorf("Error getting channel client: %v", err)
	}

	var targets []fab.Peer
	if len(Config().PeerURL) > 0 || len(Config().OrgIDs) > 0 {
		targets = a.Peers()
	}

	// executor := executor.NewConcurrent("Query Chaincode", Config().Concurrency)
	// executor.Start()
	// defer executor.Stop(true)

	// verbose := Config().Verbose || Config().Iterations == 1

	// var mutex sync.RWMutex
	// var tasks []*QueryTask
	// var errs []error
	// var wg sync.WaitGroup
	// var taskID int
	// var success int
	// for i := 0; i < Config().Iterations; i++ {
	// 	for _, args := range argsArray {
	// 		argsStruct := transform(&args)
	// 		taskID++
	// 		task := NewQueryTask(
	// 			strconv.Itoa(taskID), channelClient, targets,
	// 			a.ChaincodeID, argsStruct,
	// 			a.Printer(), verbose, Config().PrintPayloadOnly,
	// 			func(err error) {
	// 				defer wg.Done()
	// 				mutex.Lock()
	// 				if err != nil {
	// 					errs = append(errs, err)
	// 				} else {
	// 					success++
	// 				}
	// 				mutex.Unlock()
	// 			})
	// 		tasks = append(tasks, task)
	// 	}
	// }

	// numInvocations := len(tasks)
	// wg.Add(numInvocations)

	// done := make(chan bool)
	// go func() {
	// 	ticker := time.NewTicker(3 * time.Second)
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			mutex.RLock()
	// 			if len(errs) > 0 {
	// 				fmt.Printf("*** %d failed query(s) out of %d\n", len(errs), numInvocations)
	// 			}
	// 			fmt.Printf("*** %d successfull query(s) out of %d\n", success, numInvocations)
	// 			mutex.RUnlock()
	// 		case <-done:
	// 			return
	// 		}
	// 	}
	// }()

	// startTime := time.Now()

	// for _, task := range tasks {
	// 	if err := executor.Submit(task); err != nil {
	// 		return errors.Errorf("error submitting task: %s", err)
	// 	}
	// }

	// // Wait for all tasks to complete
	// wg.Wait()
	// done <- true
	// duration := time.Now().Sub(startTime)

	// if len(errs) > 0 {
	// 	fmt.Printf("\n*** %d errors querying chaincode:\n", len(errs))
	// 	for _, err := range errs {
	// 		fmt.Printf("%s\n", err)
	// 	}
	// }

	// if numInvocations > 1 {
	// 	fmt.Printf("\n")
	// 	fmt.Printf("*** ---------- Summary: ----------\n")
	// 	fmt.Printf("***   - Queries:     %d\n", numInvocations)
	// 	fmt.Printf("***   - Concurrency: %d\n", Config().Concurrency)
	// 	fmt.Printf("***   - Successfull: %d\n", success)
	// 	fmt.Printf("***   - Duration:    %s\n", duration)
	// 	fmt.Printf("***   - Rate:        %2.2f/s\n", float64(numInvocations)/duration.Seconds())
	// 	fmt.Printf("*** ------------------------------\n")
	// }

	// return nil

	var opts []channel.RequestOption
	if len(targets) > 0 {
		opts = append(opts, channel.WithTargets(targets...))
	}
	var resp []string
	for _, args := range argsArray {
		// 		argsStruct := transform(&args)
		if response, err := channelClient.Query(
			channel.Request{
				ChaincodeID: a.ChaincodeID,
				Fcn:         args.Func,
				Args:        utils.AsBytes(args.Args),
			},
			opts...,
		); err != nil {
			log.Debugf("(%s) - Error querying chaincode: %s\n", a.ChaincodeID, err)

		} else {
			resp = append(resp, string(response.Payload))
			log.Debugf("(%s) - Chaincode query was successful\n", a.ChaincodeID)
		}
	}
	if len(resp) > 0 {
		return resp[0], nil
	}
	return "", nil
}
