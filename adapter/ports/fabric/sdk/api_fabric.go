package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	ethsdk "github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	"github.com/QOSGroup/cassini/log"
	"github.com/pkg/errors"
)

const (
	errUnsuportedToken = `{"code": 404, "message": "unsupported network and token"}`
)

// ChaincodeInvoke invoke chaincode
func ChaincodeInvoke(channelID, chaincodeID string, argsArray []Args) (result string, err error) {
	log.Info("chaincode invoke...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	action, err := newChaincodeInvokeAction()
	if err != nil {
		log.Errorf("Error while initializing invokeAction: %v", err)
		return
	}

	defer action.Terminate()

	result, err = action.invoke(channelID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while calling action.invoke(): %v", err)
	}
	return
}

// ChaincodeInvokeByString call chaincode invoke of hyperledger fabric
func ChaincodeInvokeByString(channelID, chaincodeID, argsStr string) string {
	result := &CallResult{
		Code: http.StatusOK, Message: "OK"}
	argsArray, err := ArgsArray(argsStr)
	if err == nil {
		for _, args := range argsArray {
			result.Message = fmt.Sprintf("OK; %s:%s:%s", channelID, chaincodeID, args.Func)
		}
		ret, err := ChaincodeInvoke(channelID, chaincodeID, argsArray)
		if err != nil {
			result.Code = http.StatusInternalServerError
			result.Message = fmt.Sprintf("chaincode invoke error: %v", err)
		} else {
			log.Info("chaincode invoke result: ", ret)
		}
	} else {
		result.Code = http.StatusInternalServerError
		result.Message = fmt.Sprintf("args JSON parsing err: %v", err)
	}
	bytes, err := json.Marshal(result)
	if err == nil {
		log.Info(string(bytes))
		return string(bytes)
	}
	log.Errorf("%s %v", DefaultResultJSON, err)
	return DefaultResultJSON
}

// ChaincodeQuery call chaincode query of hyperledger fabric
func ChaincodeQuery(channelID, chaincodeID string, argsArray []Args) (result string, err error) {
	log.Info("chaincode query...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	var action *queryAction
	action, err = newQueryAction()
	if err != nil {
		log.Errorf("Error while initializing queryAction: %v", err)
		return
	}

	defer action.Terminate()

	result, err = action.query(channelID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while running queryAction: %v", err)
	} else if result == "" {
		err = errors.New("transaction not found")
	}
	if strings.HasPrefix(result, "v") {
		if i := strings.Index(result, ":"); i > -1 {
			result = result[i+1:]
		}
	}
	return
}

// ChaincodeQueryByString call chaincode query of hyperledger fabric
func ChaincodeQueryByString(channelID, chaincodeID, argsStr string) string {
	log.Infof("%s; %s; %s", channelID, chaincodeID, argsStr)
	result := &CallResult{
		Code: http.StatusOK, Message: "OK"}
	argsArray, err := ArgsArray(argsStr)
	if err == nil {
		for _, args := range argsArray {
			result.Message = fmt.Sprintf("OK; %s:%s:%s", channelID, chaincodeID, args.Func)
			log.Infof("result.Message: ", result.Message)
		}
		ret, err := ChaincodeQuery(channelID, chaincodeID, argsArray)
		if err != nil {
			result.Code = http.StatusInternalServerError
			result.Message = fmt.Sprintf("chaincode query error: %v", err)
		} else {
			result.Result = ret
			log.Info("chaincode query result: ", ret)
		}
	} else {
		result.Code = http.StatusInternalServerError
		result.Message = fmt.Sprintf("args JSON parsing err: %v", err)
	}
	bytes, err := json.Marshal(result)
	if err == nil {
		log.Info(string(bytes))
		return string(bytes)
	}
	log.Errorf("%s %v", DefaultResultJSON, err)
	return DefaultResultJSON
}

// NewAccountByString create a new account
func NewAccountByString(accountID, key, chain, token string) string {
	if strings.EqualFold(chain, "ethereum") {
		if strings.EqualFold(token, "eth") {
			return ethNewAccount(accountID, key, chain, token)
		}
	}
	return errUnsuportedToken
}

func ethNewAccount(accountID, key, chain, token string) string {
	account, err := ethsdk.NewAccount(accountID, key)
	if err != nil {
		log.Errorf("new account error: %v", err)
		return errUnsuportedToken
	}
	var strs []string
	height, err := ethsdk.EthBlockNumber()
	if err != nil {
		log.Errorf("eth get block number error: %v", err)
		return errUnsuportedToken
	}
	strs = append(strs, accountID,
		account.WalletAddress, "", chain, token, height)
	a := Args{Func: "register", Args: strs}
	var args []Args
	args = append(args, a)
	ret, err := ChaincodeInvoke(Config().ChannelID, "wallet", args)
	if err != nil {
		log.Errorf("new account error: %v", err)
		return errUnsuportedToken
	}
	return ret
}
