package walletextension

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/gorilla/websocket"
)

// ObscuroFacade is a server that inverts the encryption and decryption performed by WalletExtension, so that the forwarded
// Ethereum JSON-RPC requests can be understood by a regular Geth node.
type ObscuroFacade struct {
	enclaveEciesPrivateKey *ecies.PrivateKey
	gethWebsocketAddr      string
	viewingKeyChannel      <-chan ViewingKey
	viewingKeyEcies        *ecies.PublicKey
	viewingKeyAddr         string // The address the public key that signed the viewing key corresponds to.
	upgrader               websocket.Upgrader
	server                 *http.Server
}

func NewObscuroFacade(
	enclavePrivateKey *ecdsa.PrivateKey,
	gethWebsocketAddr string,
	viewingKeyChannel <-chan ViewingKey,
) *ObscuroFacade {
	enclaveEciesPrivateKey := ecies.ImportECDSA(enclavePrivateKey)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &ObscuroFacade{
		enclaveEciesPrivateKey: enclaveEciesPrivateKey,
		gethWebsocketAddr:      gethWebsocketAddr,
		viewingKeyChannel:      viewingKeyChannel,
		upgrader:               upgrader,
	}
}

// Serve listens for and serves Ethereum JSON-RPC requests from the wallet extension.
func (of *ObscuroFacade) Serve(hostAndPort string) {
	// We listen for the account viewing key.
	go func() {
		viewingKey := <-of.viewingKeyChannel

		// We recalculate the message signed by MetaMask.
		viewingKeyBytes := crypto.CompressPubkey(viewingKey.publicKey)
		msgToSign := signedMsgPrefix + hex.EncodeToString(viewingKeyBytes)

		// We recover the key based on the signed message and the signature.
		recoveredKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), viewingKey.signature)
		if err != nil {
			panic(fmt.Sprintf("was sent viewing key but could not validate its signature: %s", err))
		}

		of.viewingKeyAddr = crypto.PubkeyToAddress(*recoveredKey).String()
		of.viewingKeyEcies = ecies.ImportECDSAPublic(viewingKey.publicKey)
	}()

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pathRoot, of.handleWSEthJSON)
	of.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err := of.server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(fmt.Sprintf("could not serve the Obscuro facade: %s", err))
	}
}

func (of *ObscuroFacade) Shutdown() {
	if of.server != nil {
		err := of.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down Obscuro facade: %s", err)
		}
	}
}

func (of *ObscuroFacade) handleWSEthJSON(resp http.ResponseWriter, req *http.Request) {
	// TODO - Maintain a single connection over time, rather than recreating one for each request.
	connection, err := of.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not upgrade connection to a websocket connection: %s", err), httpCodeErr)
		return
	}
	defer connection.Close()

	// We read the message from the wallet extension.
	_, encryptedMessage, err := connection.ReadMessage()
	if err != nil {
		sendErr(connection, fmt.Sprintf("could not read Ethereum JSON-RPC request: %s", err))
		return
	}

	// We decrypt the JSON with the enclave's private key.
	message, err := of.enclaveEciesPrivateKey.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		sendErr(connection, fmt.Sprintf("could not decrypt Ethereum JSON-RPC request with enclave public key: %s", err))
		return
	}

	// We forward the message to the Geth node.
	gethResp, err := forwardMsgOverWebsocket(of.gethWebsocketAddr, message)
	if err != nil {
		sendErr(connection, fmt.Sprintf("could not forward request to Geth node: %s", err))
		return
	}

	// We unmarshall the JSON request to inspect it.
	var reqJSONMap map[string]interface{}
	err = json.Unmarshal(message, &reqJSONMap)
	if err != nil {
		sendErr(connection, fmt.Sprintf("could not unmarshall Ethereum JSON-RPC request to JSON: %s", err))
		return
	}

	// We encrypt the response if needed.
	// TODO - Return error gracefully instead of panicking if JSON keys are missing.
	method := reqJSONMap[reqJSONKeyMethod]
	if method == reqJSONMethodGetBalance || method == reqJSONMethodCall {
		firstParam := reqJSONMap[reqJSONKeyParams].([]interface{})[0]
		var requesterAddr string
		if method == reqJSONMethodGetBalance {
			requesterAddr = firstParam.(string)
		} else {
			requesterAddr = firstParam.(map[string]interface{})[reqJSONKeyFrom].(string)
		}

		if of.viewingKeyEcies == nil || requesterAddr != of.viewingKeyAddr {
			msg := fmt.Sprintf("enclave could not respond securely to %s request because there is no viewing key for the account", method)
			sendErr(connection, msg)
			return
		}

		gethResp, err = ecies.Encrypt(rand.Reader, of.viewingKeyEcies, gethResp, nil, nil)
		if err != nil {
			sendErr(connection, fmt.Sprintf("could not encrypt Ethereum JSON-RPC response with viewing key: %s", err))
			return
		}
	}

	// We write the message back to the wallet extension.
	err = connection.WriteMessage(websocket.TextMessage, gethResp)
	if err != nil {
		fmt.Printf("could not write JSON-RPC response: %s\n", err)
	}
}

// Sends the error message as a websocket error.
func sendErr(connection *websocket.Conn, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		respJSONKeyErr: map[string]string{
			respJSONKeyMsg: msg,
		},
	})
	if err != nil {
		panic(err)
	}

	_ = connection.WriteMessage(websocket.TextMessage, resp)
}
