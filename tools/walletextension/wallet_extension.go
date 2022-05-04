package walletextension

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"net/http"
	"strings"
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
	obscuroFacadeAddr string
	// TODO - Support multiple viewing keys. This will require the enclave to attach metadata on encrypted results
	//  to indicate which viewing key they were encrypted with.
	viewingKeyPrivate *ecdsa.PrivateKey
	// TODO - Replace this channel with port-based communication with the enclave.
	viewingKeyChannel chan<- ViewingKey
}

func NewWalletExtension(
	enclavePrivateKey *ecdsa.PrivateKey,
	obscuroFacadeAddr string,
	viewingKeyChannel chan<- ViewingKey,
) *WalletExtension {
	return &WalletExtension{
		enclavePrivateKey: enclavePrivateKey,
		obscuroFacadeAddr: obscuroFacadeAddr,
		viewingKeyChannel: viewingKeyChannel}
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc(pathRoot, we.handleHttpEthJson)

	// Handles the management of viewing keys.
	serveMux.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.Dir(staticDir))))
	serveMux.HandleFunc(pathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMux.HandleFunc(pathSubmitViewingKey, we.handleSubmitViewingKey)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

// Encrypts Ethereum JSON-RPC request, forwards it to the Geth node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleHttpEthJson(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not read JSON-RPC request body: %v\n", err), httpCodeErr)
		return
	}

	// We unmarshall the JSON request.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(body, &reqJsonMap)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not unmarshall JSON-RPC request body to JSON: %v\n", err), httpCodeErr)
		return
	}
	method := reqJsonMap[reqJsonKeyMethod]
	fmt.Println(fmt.Sprintf("Received request from wallet: %s", body))

	// We encrypt the JSON with the enclave's public key.
	fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	eciesPublicKey := ecies.ImportECDSAPublic(&we.enclavePrivateKey.PublicKey)
	encryptedBody, err := ecies.Encrypt(rand.Reader, eciesPublicKey, body, nil, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not encrypt request with enclave public key: %v\n", err), httpCodeErr)
		return
	}

	// We forward the request on to the Geth node.
	gethResp := forwardMsgOverWebsocket("ws://"+we.obscuroFacadeAddr, encryptedBody)
	// TODO - Improve error detection. We are just matching on the error message here.
	if strings.HasPrefix(string(gethResp), "enclave could not respond securely") {
		fmt.Println(string(gethResp))
		http.Error(resp, string(gethResp), httpCodeErr)
		return
	}

	// We decrypt the response if it's encrypted.
	if method == reqJsonMethodGetBalance || method == reqJsonMethodGetStorageAt {
		fmt.Println(fmt.Sprintf("🔐 Decrypting %s response from Geth node with viewing key.", method))
		eciesPrivateKey := ecies.ImportECDSA(we.viewingKeyPrivate)
		gethResp, err = eciesPrivateKey.Decrypt(gethResp, nil, nil)
		if err != nil {
			http.Error(resp, fmt.Sprintf("could not decrypt enclave response with viewing key: %v\n", err), httpCodeErr)
			return
		}
	}

	// We unmarshall the JSON response.
	var respJsonMap map[string]interface{}
	err = json.Unmarshal(gethResp, &respJsonMap)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not unmarshall enclave response to JSON: %v\n", err), httpCodeErr)
		return
	}
	fmt.Println(fmt.Sprintf("Received response from Geth node: %s", strings.TrimSpace(string(gethResp))))

	// We write the response to the client.
	_, err = resp.Write(gethResp)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not write JSON-RPC response: %v\n", err), httpCodeErr)
		return
	}
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(resp http.ResponseWriter, _ *http.Request) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not generate new keypair: %v\n", err), httpCodeErr)
		return
	}
	we.viewingKeyPrivate = viewingKeyPrivate

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	_, err = resp.Write([]byte(viewingKeyHex))
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not return viewing key public key hex to client: %v\n", err), httpCodeErr)
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not read viewing key and signed bytes from client: %v\n", err), httpCodeErr)
		return
	}

	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(body, &reqJsonMap)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not unmarshall viewing key and signed bytes from client to JSON: %v\n", err), httpCodeErr)
		return
	}
	signedBytes := []byte(reqJsonMap["signedBytes"].(string))

	viewingKey := ViewingKey{viewingKeyPublic: &we.viewingKeyPrivate.PublicKey, signedBytes: signedBytes}
	we.viewingKeyChannel <- viewingKey
}
