package sdk

import (
	"time"

	"github.com/QOSGroup/cassini/log"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/chaincode/invokeerror"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/chaincode/utils"
)

type chaincodeInvokeAction struct {
	Action
	numInvoked uint32
	done       chan bool
}

func newChaincodeInvokeAction() (*chaincodeInvokeAction, error) {
	action := &chaincodeInvokeAction{done: make(chan bool)}
	err := action.Initialize()
	return action, err
}

func (a *chaincodeInvokeAction) invoke(channelID, chaincodeID string, argsArray []Args) (string, error) {
	a.Set(channelID, chaincodeID, argsArray)
	channelClient, err := a.ChannelClient()
	if err != nil {
		return "", errors.Errorf("Error getting channel client: %v", err)
	}

	var targets []fab.Peer
	if len(Config().PeerURL) > 0 || len(Config().OrgIDs) > 0 {
		targets = a.Peers()
	}

	// var wg sync.WaitGroup
	// var mutex sync.RWMutex
	// var tasks []*invoketask.Task
	// var taskID int
	// for i := 0; i < Config().Iterations; i++ {
	// 	for _, args := range a.Args {
	// 		argsStruct := transform(&args)
	// 		taskID++
	// 		var startTime time.Time
	// 		task := invoketask.New(
	// 			strconv.Itoa(taskID), channelClient, targets,
	// 			a.ChaincodeID, argsStruct, executor,
	// 			retry.Opts{
	// 				Attempts:       Config().MaxAttempts,
	// 				InitialBackoff: Config().InitialBackoff,
	// 				MaxBackoff:     Config().MaxBackoff,
	// 				BackoffFactor:  Config().BackoffFactor,
	// 				RetryableCodes: retry.ChannelClientRetryableCodes,
	// 			},
	// 			Config().Verbose || Config().Iterations == 1,
	// 			Config().PrintPayloadOnly, a.Printer(),

	// 			func() {
	// 				startTime = time.Now()
	// 			},
	// 			func(err error) {
	// 				duration := time.Since(startTime)
	// 				defer wg.Done()
	// 				mutex.Lock()
	// 				defer mutex.Unlock()
	// 				if err != nil {
	// 					errs = append(errs, err)
	// 					failDurations = append(failDurations, duration)
	// 				} else {
	// 					success++
	// 					successDurations = append(successDurations, duration)
	// 				}
	// 			})
	// 		tasks = append(tasks, task)
	// 	}
	// }

	// numInvocations := len(tasks)

	// wg.Add(numInvocations)

	// done := make(chan bool)
	// go func() {
	// 	ticker := time.NewTicker(10 * time.Second)
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			mutex.RLock()
	// 			if len(errs) > 0 {
	// 				fmt.Printf("*** %d failed invocation(s) out of %d\n", len(errs), numInvocations)
	// 			}
	// 			fmt.Printf("*** %d successfull invocation(s) out of %d\n", success, numInvocations)
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

	// var allErrs []error
	// var attempts int
	// for _, task := range tasks {
	// 	attempts = attempts + task.Attempts()
	// 	if task.LastError() != nil {
	// 		allErrs = append(allErrs, task.LastError())
	// 	}
	// }

	// if len(errs) > 0 {
	// 	fmt.Printf("\n*** %d errors invoking chaincode:\n", len(errs))
	// 	for _, err := range errs {
	// 		fmt.Printf("%s\n", err)
	// 	}
	// } else if len(allErrs) > 0 {
	// 	fmt.Printf("\n*** %d transient errors invoking chaincode:\n", len(allErrs))
	// 	for _, err := range allErrs {
	// 		fmt.Printf("%s\n", err)
	// 	}
	// }

	// if numInvocations > 1 {
	// 	fmt.Printf("\n")
	// 	fmt.Printf("*** ---------- Summary: ----------\n")
	// 	fmt.Printf("***   - Invocations:     %d\n", numInvocations)
	// 	fmt.Printf("***   - Concurrency:     %d\n", Config().Concurrency)
	// 	fmt.Printf("***   - Successfull:     %d\n", success)
	// 	fmt.Printf("***   - Total attempts:  %d\n", attempts)
	// 	fmt.Printf("***   - Duration:        %2.2fs\n", duration.Seconds())
	// 	fmt.Printf("***   - Rate:            %2.2f/s\n", float64(numInvocations)/duration.Seconds())
	// 	fmt.Printf("***   - Average:         %2.2fs\n", average(append(successDurations, failDurations...)))
	// 	fmt.Printf("***   - Average Success: %2.2fs\n", average(successDurations))
	// 	fmt.Printf("***   - Average Fail:    %2.2fs\n", average(failDurations))
	// 	fmt.Printf("***   - Min Success:     %2.2fs\n", min(successDurations))
	// 	fmt.Printf("***   - Max Success:     %2.2fs\n", max(successDurations))
	// 	fmt.Printf("*** ------------------------------\n")
	// }

	var opts []channel.RequestOption
	// opts = append(opts, channel.WithRetry(t.retryOpts))
	// opts = append(opts, channel.WithBeforeRetry(func(err error) {
	// 	t.attempt++
	// }))
	if len(targets) > 0 {
		opts = append(opts, channel.WithTargets(targets...))
	}
	var resp []string
	for _, args := range argsArray {
		response, err := channelClient.Execute(
			channel.Request{
				ChaincodeID: chaincodeID,
				Fcn:         args.Func,
				Args:        utils.AsBytes(args.Args),
			},
			opts...,
		)
		if err != nil {
			return "", invokeerror.Errorf(invokeerror.TransientError, "SendTransactionProposal return error: %v", err)
		}

		txID := string(response.TransactionID)

		switch pb.TxValidationCode(response.TxValidationCode) {
		case pb.TxValidationCode_VALID:
			resp = append(resp, string(response.Payload))
			log.Infof("(%s) - Successfully committed transaction [%s] ...\n", txID, response.TransactionID)
		case pb.TxValidationCode_DUPLICATE_TXID, pb.TxValidationCode_MVCC_READ_CONFLICT, pb.TxValidationCode_PHANTOM_READ_CONFLICT:
			log.Infof("(%s) - Transaction commit failed for [%s] with code [%s]. This is most likely a transient error.\n", txID, response.TransactionID, response.TxValidationCode)
			return "", invokeerror.Wrapf(invokeerror.TransientError, errors.New("Duplicate TxID"), "invoke Error received from eventhub for TxID [%s]. Code: %s", response.TransactionID, response.TxValidationCode)
		default:
			log.Infof("(%s) - Transaction commit failed for [%s] with code [%s].\n", txID, response.TransactionID, response.TxValidationCode)
			return "", invokeerror.Wrapf(invokeerror.PersistentError, errors.New("error"), "invoke Error received from eventhub for TxID [%s]. Code: %s", response.TransactionID, response.TxValidationCode)
		}
	}
	if len(resp) > 0 {
		return resp[0], nil
	}
	return "", nil
}

func average(durations []time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}

	var total float64
	for _, duration := range durations {
		total += duration.Seconds()
	}
	return total / float64(len(durations))
}

func min(durations []time.Duration) float64 {
	min, _ := minMax(durations)
	return min
}

func max(durations []time.Duration) float64 {
	_, max := minMax(durations)
	return max
}

func minMax(durations []time.Duration) (min float64, max float64) {
	for _, duration := range durations {
		if min == 0 || min > duration.Seconds() {
			min = duration.Seconds()
		}
		if max == 0 || max < duration.Seconds() {
			max = duration.Seconds()
		}
	}
	return
}
