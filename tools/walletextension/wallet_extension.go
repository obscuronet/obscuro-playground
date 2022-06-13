package walletextension

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/obscuroclient"

	"github.com/gorilla/websocket"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	pathRoot               = "/"
	PathReady              = "/ready/"
	pathViewingKeys        = "/viewingkeys/"
	PathGenerateViewingKey = "/generateviewingkey/"
	PathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "./tools/walletextension/static"

	reqJSONKeyMethod        = "method"
	ReqJSONMethodGetBalance = "eth_getBalance"
	ReqJSONMethodCall       = "eth_call"
	respJSONKeyErr          = "error"
	respJSONKeyMsg          = "message"
	RespJSONKeyResult       = "result"
	httpCodeErr             = 500

	Localhost         = "127.0.0.1"
	websocketProtocol = "ws://"

	// EnclavePublicKeyHex is the public key of the enclave.
	// TODO - Retrieve this key from the management contract instead.
	EnclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

// TODO - Display error in browser if Metamask is not enabled (i.e. `ethereum` object is not available in-browser).

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePublicKey *ecdsa.PublicKey // The public key used to encrypt requests for the enclave.
	hostAddr         string           // The address on which the Obscuro host can be reached.
	hostClient       obscuroclient.Client
	// TODO - Support multiple viewing keys. This will require the enclave to attach metadata on encrypted results
	//  to indicate which viewing key they were encrypted with.
	viewingPublicKeyBytes  []byte
	viewingPrivateKeyEcies *ecies.PrivateKey
	server                 *http.Server
}

func NewWalletExtension(config Config) *WalletExtension {
	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	return &WalletExtension{
		enclavePublicKey: &enclavePrivateKey.PublicKey,
		hostAddr:         config.NodeRPCWebsocketAddress,
		hostClient:       obscuroclient.NewClient(config.NodeRPCHTTPAddress),
	}
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()
	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc(pathRoot, we.handleHTTPEthJSON)
	serveMux.HandleFunc(PathReady, we.handleReady)
	// Handles the management of viewing keys.
	serveMux.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.Dir(staticDir))))
	serveMux.HandleFunc(PathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMux.HandleFunc(PathSubmitViewingKey, we.handleSubmitViewingKey)
	we.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err := we.server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (we *WalletExtension) Shutdown() {
	if we.server != nil {
		err := we.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down wallet extension: %s", err)
		}
	}
}

// Used to check whether the server is ready.
func (we *WalletExtension) handleReady(http.ResponseWriter, *http.Request) {}

// Encrypts Ethereum JSON-RPC request, forwards it to the Obscuro node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleHTTPEthJSON(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read JSON-RPC request body: %s", err))
		return
	}

	// We unmarshall the JSON request.
	var reqJSONMap map[string]interface{}
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall JSON-RPC request body to JSON: %s", err))
		return
	}
	method := reqJSONMap[reqJSONKeyMethod]
	fmt.Printf("Received request from wallet: %s\n", body)

	// TODO - Reenable encryption of requests.
	//// We encrypt the JSON with the enclave's public key.
	//fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	//eciesPublicKey := ecies.ImportECDSAPublic(we.enclavePublicKey)
	//encryptedBody, err := ecies.Encrypt(rand.Reader, eciesPublicKey, body, nil, nil)
	//if err != nil {
	//	logAndSendErr(resp, fmt.Sprintf("could not encrypt request with enclave public key: %s", err))
	//	return
	//}

	// We forward the request on to the Obscuro node.
	nodeResp, err := forwardMsgOverWebsocket(websocketProtocol+we.hostAddr, body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("received error response when forwarding request to node at %s: %s", we.hostAddr, err))
		return
	}

	// We unmarshall the JSON response.
	var respJSONMap map[string]interface{}
	err = json.Unmarshal(nodeResp, &respJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall enclave response to JSON: %s", err))
		return
	}

	// We report any errors from the request.
	if respJSONMap[respJSONKeyErr] != nil {
		logAndSendErr(resp, respJSONMap[respJSONKeyErr].(map[string]interface{})[respJSONKeyMsg].(string))
		return
	}

	// We decrypt the result field if it's encrypted.
	if method == ReqJSONMethodGetBalance || method == ReqJSONMethodCall {
		fmt.Printf("🔐 Decrypting %s response from Obscuro node with viewing key.\n", method)

		encryptedResult := common.Hex2Bytes(respJSONMap[RespJSONKeyResult].(string))
		decryptedResult, err := we.viewingPrivateKeyEcies.Decrypt(encryptedResult, nil, nil)
		if err != nil {
			logAndSendErr(resp, fmt.Sprintf("could not decrypt enclave response with viewing key: %s", err))
			return
		}

		respJSONMap[RespJSONKeyResult] = string(decryptedResult)
	}

	clientResponse, err := json.Marshal(respJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not marshal JSON response to present to the client: %s", err))
		return
	}
	fmt.Printf("Received response from Obscuro node: %s\n", strings.TrimSpace(string(clientResponse)))

	// We write the response to the client.
	_, err = resp.Write(clientResponse)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not write JSON-RPC response: %s", err))
		return
	}
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(resp http.ResponseWriter, _ *http.Request) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}
	we.viewingPublicKeyBytes = crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	we.viewingPrivateKeyEcies = ecies.ImportECDSA(viewingKeyPrivate)

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	_, err = resp.Write([]byte(viewingKeyHex))
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return viewing key public key hex to client: %s", err))
		return
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read viewing key and signature from client: %s", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall viewing key and signature from client to JSON: %s", err))
		return
	}

	// We have to drop the leading "0x", and transform the V from 27/28 to 0/1.
	signature, err := hex.DecodeString(reqJSONMap["signature"][2:])
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}
	signature[64] -= 27

	var rpcErr error
	err = we.hostClient.Call(&rpcErr, obscuroclient.RPCAddViewingKey, we.viewingPublicKeyBytes, signature)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not add viewing key: %s", err))
		return
	}
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionPort     int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
}

func forwardMsgOverWebsocket(url string, msg []byte) ([]byte, error) {
	connection, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	defer resp.Body.Close()

	err = connection.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return nil, err
	}

	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}
