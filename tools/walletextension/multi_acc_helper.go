package walletextension

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

const (
	reqJSONKeyFrom      = "from"
	reqJSONKeyData      = "data"
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"
)

// multi_acc_helper provides a single location for code that helps wallet extension in determining the appropriate account
//	to use to send a request when multiple are registered

// suggestAccountClient works through various methods to try and guess which available client to use for a request, returns nil if none found
func suggestAccountClient(req *rpcRequest, accClients map[common.Address]*rpcclientlib.EncRPCClient) *rpcclientlib.EncRPCClient {
	if len(accClients) == 1 {
		for _, client := range accClients {
			// return the first (and only) client
			return client
		}
	}

	paramsMap, err := parseParams(req.params)
	if err != nil {
		// no further info to deduce calling client
		return nil
	}

	// check if request params had a "from" address and if we had a client for that address
	fromClient, found := checkForFromField(paramsMap, accClients)
	if found {
		return fromClient
	}

	if req.method == rpcclientlib.RPCCall {
		// Otherwise, we search the `data` field for an address matching a registered viewing key.
		addr, err := searchDataFieldForAccount(paramsMap, accClients)
		if err == nil {
			return accClients[*addr]
		}
	}

	// todo: add other mechanisms for determining the correct account to use. E.g. we may want to start caching and
	// 	 	recent transaction hashes for accounts so that receipt lookups know which acc to use

	return nil
}

func checkForFromField(paramsMap map[string]interface{}, accClients map[common.Address]*rpcclientlib.EncRPCClient) (*rpcclientlib.EncRPCClient, bool) {
	fromVal, found := paramsMap[reqJSONKeyFrom]
	if !found {
		return nil, false
	}

	fromStr, ok := fromVal.(string)
	if !ok {
		return nil, false
	}

	fromAddr := common.HexToAddress(fromStr)
	client, found := accClients[fromAddr]
	return client, found
}

// Extracts the arguments from the request's `data` field. If any of them, after removing padding, match the viewing
// key address, we return that address. Otherwise, we return nil.
func searchDataFieldForAccount(callParams map[string]interface{}, accClients map[common.Address]*rpcclientlib.EncRPCClient) (*common.Address, error) {
	// We ensure that the `data` field is present.
	data := callParams[reqJSONKeyData]
	if data == nil {
		return nil, fmt.Errorf("eth_call request did not have its `data` field set")
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("eth_call request's `data` field was not of the expected type `string`")
	}

	// We check that the data field is long enough before removing the leading "0x" (1 bytes/2 chars) and the method ID
	// (4 bytes/8 chars).
	if len(dataString) < 10 {
		return nil, nil //nolint:nilnil
	}
	dataString = dataString[10:]

	// We split up the arguments in the `data` field.
	var dataArgs []string
	for i := 0; i < len(dataString); i += ethCallPaddedArgLen {
		if i+ethCallPaddedArgLen > len(dataString) {
			break
		}
		dataArgs = append(dataArgs, dataString[i:i+ethCallPaddedArgLen])
	}

	// We iterate over the arguments, looking for an argument that matches a viewing key address
	for _, dataArg := range dataArgs {
		// If the argument doesn't have the correct padding, it's not an address.
		if !strings.HasPrefix(dataArg, ethCallAddrPadding) {
			continue
		}

		maybeAddress := common.HexToAddress(dataArg[len(ethCallAddrPadding):])
		if _, ok := accClients[maybeAddress]; ok {
			return &maybeAddress, nil
		}
	}

	return nil, fmt.Errorf("no known account found in data bytes")
}

// Many eth RPC requests provide params as first argument in a json map with similar fields (e.g. a `from` field)
func parseParams(args []interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no params found to unmarshal")
	}

	// only interested in trying first arg
	params, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("first arg was not a byte array")
		}

		err := json.Unmarshal(callParamsJSON, &params)
		if err != nil {
			return nil, fmt.Errorf("first arg couldn't be unmarshalled into a params map")
		}
	}

	return params, nil
}
