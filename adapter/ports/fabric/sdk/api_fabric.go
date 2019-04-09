package sdk

import (
	"encoding/json"
	"fmt"
	"strings"

	ethsdk "github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	"github.com/QOSGroup/cassini/log"
	"github.com/pkg/errors"
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
	log.Infof("chaincode invoke: %s; %s; %s", channelID, chaincodeID, argsStr)
	args, err := ArgsArray(argsStr)
	if err == nil {
		var ret string
		ret, err = ChaincodeInvoke(channelID, chaincodeID, args)
		if err == nil {
			log.Info("chaincode invoke result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
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
	log.Infof("chaincode query: %s; %s; %s", channelID, chaincodeID, argsStr)
	argsArray, err := ArgsArray(argsStr)
	if err == nil {
		var ret string
		ret, err = ChaincodeQuery(channelID, chaincodeID, argsArray)
		if err == nil {
			log.Info("chaincode query result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
}

// RegisterWalletByString create a new account
func RegisterWalletByString(accountID, key, chain, token string) string {
	if strings.EqualFold(chain, "ethereum") {
		if strings.EqualFold(token, "eth") {
			return ethNewAccount(accountID, key, chain, token)
		}
	}
	return errUnsuportedToken
}

// ImportTokenByString query token info from chain
func ImportTokenByString(chain, tokenAddress string) string {
	token, err := ethsdk.ImportToken(chain, tokenAddress)
	if err != nil {
		log.Errorf("import token error: %v", err)
		return errUnsuportedToken
	}
	ret := CallResult{Code: 200, Message: "OK", Result: token}
	bytes, err := json.Marshal(&ret)
	if err != nil {
		log.Errorf("import token error: %v", err)
		return errUnsuportedToken
	}
	return string(bytes)
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
		account.WalletAddress, chain, token, height)
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
