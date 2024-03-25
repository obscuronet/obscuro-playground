package httpapi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/node"
	"github.com/ten-protocol/go-ten/tools/walletextension/rpcapi"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common/httputil"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// NewHTTPRoutes returns the http specific routes
// todo - move these to the rpc framework.
func NewHTTPRoutes(walletExt *rpcapi.Services) []node.Route {
	return []node.Route{
		{
			Name: common.PathReady,
			Func: httpHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathGetMessage,
			Func: httpHandler(walletExt, getMessageRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathAuthenticate,
			Func: httpHandler(walletExt, authenticateRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathQuery,
			Func: httpHandler(walletExt, queryRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathRevoke,
			Func: httpHandler(walletExt, revokeRequestHandler),
		},
		{
			Name: common.PathHealth,
			Func: httpHandler(walletExt, healthRequestHandler),
		},
		{
			Name: common.PathNetworkHealth,
			Func: httpHandler(walletExt, networkHealthRequestHandler),
		},
		{
			Name: common.PathVersion,
			Func: httpHandler(walletExt, versionRequestHandler),
		},
	}
}

func httpHandler(
	walletExt *rpcapi.Services,
	fun func(walletExt *rpcapi.Services, conn UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		httpRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for http requests
func httpRequestHandler(walletExt *rpcapi.Services, resp http.ResponseWriter, req *http.Request, fun func(walletExt *rpcapi.Services, conn UserConn)) {
	if walletExt.IsStopping() {
		return
	}
	if httputil.EnableCORS(resp, req) {
		return
	}
	userConn := NewUserConnHTTP(resp, req, walletExt.Logger())
	fun(walletExt, userConn)
}

// readyRequestHandler is used to check whether the server is ready
func readyRequestHandler(_ *rpcapi.Services, _ UserConn) {}

// This function handles request to /join endpoint. It is responsible to create new user (new key-pair) and store it to the db
func joinRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// audit()
	// todo (@ziga) add protection against DDOS attacks
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// generate new key-pair and store it in the database
	userID, err := walletExt.GenerateAndStoreNewUser()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal Error"))
		walletExt.Logger().Error("error creating new user", log.ErrKey, err)
	}

	// write hex encoded userID in the response
	err = conn.WriteResponse([]byte(hexutils.BytesToHex(userID)))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /authenticate endpoint.
// In the request we receive message, signature and address in JSON as request body and userID and address as query parameters
// We then check if message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func authenticateRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// read the request
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// get the text that was signed and signature
	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal request body - %w", err))
		return
	}

	// get signature from the request and remove leading two bytes (0x)
	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to decode signature - %w", err))
		return
	}

	// get address from the request
	address, ok := reqJSONMap[common.JSONKeyAddress]
	if !ok || address == "" {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to read address field from the request"))
		return
	}

	// get an optional type of the message that was signed
	messageTypeValue := common.DefaultGatewayAuthMessageType
	if typeFromRequest, ok := reqJSONMap[common.JSONKeyType]; ok && typeFromRequest != "" {
		messageTypeValue = typeFromRequest
	}

	// check if a message type is valid
	messageType, ok := viewingkey.SignatureTypeMap[messageTypeValue]
	if !ok {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid message type: %s", messageTypeValue))
	}

	// read userID from query params
	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("malformed query: 'u' required - representing encryption token - %w", err))
		return
	}

	// check signature and add address and signature for that user
	err = walletExt.AddAddressToUser(userID, address, signature, messageType)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error(fmt.Sprintf("error adding address: %s to user: %s with signature: %s", address, userID, signature))
		return
	}
	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /query endpoint.
// In the query parameters address and userID are required. We check if provided address is registered for given userID
// and return true/false in json response
func queryRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}
	address, err := getQueryParameter(conn.ReadRequestParams(), common.AddressQueryParameter)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("address ('a') not found in query parameters"))
		walletExt.Logger().Error("address ('a') not found in query parameters", log.ErrKey, err)
		return
	}
	// check if address length is correct
	if len(address) != common.EthereumAddressLen {
		handleError(conn, walletExt.Logger(), fmt.Errorf("provided address length is %d, expected: %d", len(address), common.EthereumAddressLen))
		return
	}

	// check if this account is registered with given user
	found, err := walletExt.UserHasAccount(userID, address)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("error during checking if account exists for user", "userID", userID, log.ErrKey, err)
	}

	// create and write the response
	res := struct {
		Status bool `json:"status"`
	}{Status: found}

	msg, err := json.Marshal(res)
	if err != nil {
		handleError(conn, walletExt.Logger(), err)
		return
	}

	err = conn.WriteResponse(msg)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /revoke endpoint.
// It requires userID as query parameter and deletes given user and all associated viewing keys
func revokeRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}

	// delete user and accounts associated with it from the database
	err = walletExt.DeleteUser(userID)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("unable to delete user", "userID", userID, log.ErrKey, err)
		return
	}

	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /health endpoint.
func healthRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	// TODO: connect to database and check if it is healthy
	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /network-health endpoint.
func networkHealthRequestHandler(walletExt *rpcapi.Services, userConn UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	healthStatus, err := walletExt.GetTenNodeHealthStatus()

	data, err := json.Marshal(map[string]interface{}{
		"result": healthStatus,
		"error":  err,
	})
	if err != nil {
		walletExt.Logger().Error("error marshaling response", log.ErrKey, err)
		return
	}

	err = userConn.WriteResponse(data)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /version endpoint.
func versionRequestHandler(walletExt *rpcapi.Services, userConn UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	err = userConn.WriteResponse([]byte(walletExt.Version()))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// getMessageRequestHandler handles request to /get-message endpoint.
func getMessageRequestHandler(walletExt *rpcapi.Services, conn UserConn) {
	// read the request
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// read body of the request
	var reqJSONMap map[string]interface{}
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal address request - %w", err))
		return
	}

	// get address from the request
	encryptionToken, ok := reqJSONMap[common.JSONKeyEncryptionToken]
	if !ok || len(encryptionToken.(string)) != common.MessageUserIDLen {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to read encryptionToken field from the request or it is not of correct length"))
		return
	}

	// get formats from the request, if present
	var formatsSlice []string
	if formatsInterface, ok := reqJSONMap[common.JSONKeyFormats]; ok {
		formats, ok := formatsInterface.([]interface{})
		if !ok {
			handleError(conn, walletExt.Logger(), fmt.Errorf("formats field is not an array"))
			return
		}

		for _, f := range formats {
			formatStr, ok := f.(string)
			if !ok {
				handleError(conn, walletExt.Logger(), fmt.Errorf("format value is not a string"))
				return
			}
			formatsSlice = append(formatsSlice, formatStr)
		}
	}

	userID := hexutils.HexToBytes(encryptionToken.(string))
	if len(userID) != viewingkey.UserIDLength {
		return
	}

	message, err := walletExt.GenerateUserMessageToSign(userID, formatsSlice)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("error getting message", log.ErrKey, err)
		return
	}

	// create the response structure
	type JSONResponse struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}

	// get string representation of the message format
	messageFormat := viewingkey.GetBestFormat(formatsSlice)
	messageFormatString := viewingkey.GetSignatureTypeString(messageFormat)

	response := JSONResponse{
		Message: message,
		Type:    messageFormatString,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		walletExt.Logger().Error("error marshaling JSON response", log.ErrKey, err)
		return
	}

	err = conn.WriteResponse(responseBytes)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}
