package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Route defines the path plus handler for a given path
type Route struct {
	Name string
	Func func(resp http.ResponseWriter, req *http.Request)
}

// NewHTTPRoutes returns the http specific routes
func NewHTTPRoutes(walletExt *walletextension.WalletExtension) []Route {
	return []Route{
		{
			Name: common.PathRoot,
			Func: httpHandler(walletExt, ethRequestHandler),
		},
		{
			Name: common.PathReady,
			Func: httpHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.PathGenerateViewingKey,
			Func: httpHandler(walletExt, generateViewingKeyRequestHandler),
		},

		{
			Name: common.PathSubmitViewingKey,
			Func: httpHandler(walletExt, submitViewingKeyRequestHandler),
		},

		{
			Name: common.PathAuthenticate,
			Func: httpHandler(walletExt, authenticateRequestHandler),
		},

		{
			Name: common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
		},
	}
}

func httpHandler(
	walletExt *walletextension.WalletExtension,
	fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		httpRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for http requests
func httpRequestHandler(walletExt *walletextension.WalletExtension, resp http.ResponseWriter, req *http.Request, fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn)) {
	if walletExt.IsStopping() {
		return
	}
	if httputil.EnableCORS(resp, req) {
		return
	}
	userConn := userconn.NewUserConnHTTP(resp, req, walletExt.Logger())
	fun(walletExt, userConn)
}

// NewWSRoutes returns the WS specific routes
func NewWSRoutes(walletExt *walletextension.WalletExtension) []Route {
	return []Route{
		{
			Name: common.PathRoot,
			Func: wsHandler(walletExt, ethRequestHandler),
		},
		{
			Name: common.PathReady,
			Func: wsHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.PathGenerateViewingKey,
			Func: wsHandler(walletExt, generateViewingKeyRequestHandler),
		},

		{
			Name: common.PathSubmitViewingKey,
			Func: wsHandler(walletExt, submitViewingKeyRequestHandler),
		},

		{
			Name: common.PathAuthenticate,
			Func: wsHandler(walletExt, authenticateRequestHandler),
		},

		{
			Name: common.PathJoin,
			Func: wsHandler(walletExt, joinRequestHandler),
		},
	}
}

func wsHandler(
	walletExt *walletextension.WalletExtension,
	fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		wsRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for WS requests
func wsRequestHandler(walletExt *walletextension.WalletExtension, resp http.ResponseWriter, req *http.Request, fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn)) {
	if walletExt.IsStopping() {
		return
	}

	userConn, err := userconn.NewUserConnWS(resp, req, walletExt.Logger())
	if err != nil {
		return
	}
	// We handle requests in a loop until the connection is closed on the client side.
	for !userConn.IsClosed() {
		fun(walletExt, userConn)
	}
}

// ethRequestHandler parses the user eth request, passes it on to the WE to proxy it and processes the response
func ethRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		return
	}

	request, err := parseRequest(body)
	if err != nil {
		conn.HandleError(err.Error())
		return
	}
	walletExt.Logger().Debug("REQUEST", "method", request.Method, "body", string(body))

	if request.Method == rpc.Subscribe && !conn.SupportsSubscriptions() {
		conn.HandleError(common.ErrSubscribeFailHTTP)
		return
	}

	// todo (@pedro) remove this conn dependency
	response, err := walletExt.ProxyEthRequest(request, conn)
	if err != nil {
		walletExt.Logger().Error("error while proxying request", log.ErrKey, err)
		response = common.CraftErrorResponse(err)
	}

	rpcResponse, err := json.Marshal(response)
	if err != nil {
		conn.HandleError(fmt.Sprintf("failed to remarshal RPC response to return to caller: %s", err))
		return
	}
	walletExt.Logger().Info(fmt.Sprintf("Forwarding %s response from Obscuro node: %s", request.Method, rpcResponse))

	err = conn.WriteResponse(rpcResponse)
	if err != nil {
		return
	}
}

// readyRequestHandler is used to check whether the server is ready
func readyRequestHandler(_ *walletextension.WalletExtension, _ userconn.UserConn) {}

// generateViewingKeyRequestHandler parses the gen vk request
func generateViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		conn.HandleError(fmt.Sprintf("could not unmarshal address request - %s", err))
		return
	}

	address := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	pubViewingKey, err := walletExt.GenerateViewingKey(address)
	if err != nil {
		conn.HandleError(fmt.Sprintf("unable to generate vieweing key: %s", err))
		return
	}

	err = conn.WriteResponse([]byte(pubViewingKey))
	if err != nil {
		return
	}
}

// submitViewingKeyRequestHandler submits the viewing key and signed bytes to the WE
func submitViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	body, err := userConn.ReadRequest()
	if err != nil {
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not unmarshal address and signature from client to JSON: %s", err))
		return
	}
	accAddress := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	err = walletExt.SubmitViewingKey(accAddress, signature)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not submit viewing key - %s", err))
		return
	}

	err = userConn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		return
	}
}

func authenticateRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// check if the text is well-formed and extract signature and message
	body, err := userConn.ReadRequest()
	if err != nil {
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}

	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	message, ok := reqJSONMap[common.JSONKeyMessage]
	if !ok || message == "" {
		userConn.HandleError("message not found in the request")
		return
	}

	// read userID from query params
	userID, err := getUser(userConn.ReadRequestParams())
	if err != nil {
		userConn.HandleError("userID not found in the request")
		return
	}

	// get userID and address from message
	messageUserID := ""
	messageAddressHex := ""
	regex := regexp.MustCompile(`^Register\s(\w+)\sfor\s(\w+)$`)
	if regex.MatchString(message) {
		params := regex.FindStringSubmatch(message)
		messageUserID = params[1]
		messageAddressHex = params[2]
	} else {
		userConn.HandleError(fmt.Sprintf("Submitted message is not in the correct format: %s", message))
	}

	// check if userID corresponds to the one in the message
	if userID != messageUserID || messageUserID == "" {
		userConn.HandleError(fmt.Sprintf("User in submitted message (%s) does not match user provided in the request (%s)", messageUserID, userID))
	}

	// get the address from signature

	// prefix the message like in the personal_sign method
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	if len(signature) != 65 {
		fmt.Println("Signature must be 65 bytes long")
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		fmt.Println(err)
	}

	addressFromSignature := crypto.PubkeyToAddress(*pubKey)
	addressFromMessage := gethcommon.HexToAddress(messageAddressHex)

	// verify that message was signed by the same address as in the message
	if addressFromSignature != addressFromMessage {
		fmt.Println("address from signature is not the same as address from message")
		return
	}

	// save the data for this specific userID

	// get privateKey for userID
	userPrivateKey, err := walletExt.Storage.GetUnauthenticatedUserPrivateKey(userID)
	if err != nil {
		fmt.Println("Error getting user private key")
	}
	if len(userPrivateKey) == 0 {
		fmt.Println("Received private key with length 0")
	}
	// store all the fields in the database
	err = walletExt.Storage.StoreAuthenticatedDataForUser(userID, userPrivateKey, addressFromSignature.Bytes(), message, string(signature))
	if err != nil {
		fmt.Println("Unable to store data for user: ", userID)
	}
}

func joinRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// todo (@ziga) add protection against DDOS attacks
	_, err := userConn.ReadRequest()
	if err != nil {
		return
	}

	// generate new key-pair
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}

	// create UserID
	// todo - is hash of public key ok to be user id? (public keys are usually shared with others and others can then get userID?)
	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	userID := crypto.Keccak256Hash(viewingPublicKeyBytes)

	// save UserID and PrivateKey to the database
	vk := &rpc.ViewingKey{
		Account:    nil,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	err = walletExt.Storage.SaveUserVK(userID.Hex(), vk, "")
	if err != nil {
		userConn.HandleError(fmt.Sprintf("failed to save user to the database: %s", err))
		return
	}

	err = userConn.WriteResponse([]byte(userID.Hex()))
	if err != nil {
		return
	}
}
