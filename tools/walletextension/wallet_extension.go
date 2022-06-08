package walletextension

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	pathViewingKeys        = "/viewingkeys/"
	pathGenerateViewingKey = "/generateviewingkey/"
	pathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "./tools/walletextension/static"
)

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePrivateKey *ecdsa.PrivateKey
	nodeAddr          string // The address on which the node (or facade) can be reached.
	// TODO - Support multiple viewing keys. This will require the enclave to attach metadata on encrypted results
	//  to indicate which viewing key they were encrypted with.
	viewingKeyPrivate      *ecdsa.PrivateKey
	viewingKeyPrivateEcies *ecies.PrivateKey
	// TODO - Replace this channel with port-based communication with the enclave.
	viewingKeyChannel chan<- ViewingKey
	server            *http.Server
}

func NewWalletExtension(
	enclavePrivateKey *ecdsa.PrivateKey,
	nodeAddr string,
	viewingKeyChannel chan<- ViewingKey,
) *WalletExtension {
	return &WalletExtension{
		enclavePrivateKey: enclavePrivateKey,
		nodeAddr:          nodeAddr,
		viewingKeyChannel: viewingKeyChannel,
	}
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()
	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc(pathRoot, we.handleHTTPEthJSON)
	// Handles the management of viewing keys.
	serveMux.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.Dir(staticDir))))
	serveMux.HandleFunc(pathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMux.HandleFunc(pathSubmitViewingKey, we.handleSubmitViewingKey)
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

// Encrypts Ethereum JSON-RPC request, forwards it to the Geth node over a websocket, and decrypts the response if needed.
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

	// We encrypt the JSON with the enclave's public key.
	fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	eciesPublicKey := ecies.ImportECDSAPublic(&we.enclavePrivateKey.PublicKey)
	encryptedBody, err := ecies.Encrypt(rand.Reader, eciesPublicKey, body, nil, nil)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not encrypt request with enclave public key: %s", err))
		return
	}

	// We forward the request on to the Geth node.
	gethResp, err := forwardMsgOverWebsocket(websocketProtocol+we.nodeAddr, encryptedBody)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("received error response when forwarding request to node at %s: %s", we.nodeAddr, err))
		return
	}

	// This is just a temporary unmarshalling. We need to unmarshall once to check if we got an error response, then
	// unmarshall again once we've decrypted the response if needed, below.
	var respJSONMapTemp map[string]interface{}
	err = json.Unmarshal(gethResp, &respJSONMapTemp)
	// A nil error indicates that this was valid JSON, and not an encrypted payload.
	if err == nil && respJSONMapTemp[respJSONKeyErr] != nil {
		logAndSendErr(resp, respJSONMapTemp[respJSONKeyErr].(map[string]interface{})[respJSONKeyMsg].(string))
		return
	}

	// We decrypt the response if it's encrypted.
	if method == reqJSONMethodGetBalance || method == reqJSONMethodCall {
		fmt.Printf("🔐 Decrypting %s response from Geth node with viewing key.\n", method)
		gethResp, err = we.viewingKeyPrivateEcies.Decrypt(gethResp, nil, nil)
		if err != nil {
			logAndSendErr(resp, fmt.Sprintf("could not decrypt enclave response with viewing key: %s", err))
			return
		}
	}

	// We unmarshall the JSON response.
	var respJSONMap map[string]interface{}
	err = json.Unmarshal(gethResp, &respJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall enclave response to JSON: %s", err))
		return
	}
	fmt.Printf("Received response from Geth node: %s\n", strings.TrimSpace(string(gethResp)))

	// We write the response to the client.
	_, err = resp.Write(gethResp)
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
	we.viewingKeyPrivate = viewingKeyPrivate
	we.viewingKeyPrivateEcies = ecies.ImportECDSA(viewingKeyPrivate)

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

	viewingKey := ViewingKey{publicKey: &we.viewingKeyPrivate.PublicKey, signature: signature}
	we.viewingKeyChannel <- viewingKey
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
