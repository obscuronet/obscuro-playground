package walletextension

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"

	"github.com/gorilla/websocket"
)

const (
	reqJSONKeyMethod        = "method"
	reqJSONMethodChainID    = "eth_chainId"
	reqJSONMethodGetBalance = "eth_getBalance"
	reqJSONMethodCall       = "eth_call"
	reqJSONKeyParams        = "params"
	reqJSONKeyTo            = "to"
	reqJSONKeyFrom          = "from"
	respJSONKeyErr          = "error"
	respJSONKeyMsg          = "message"
	pathRoot                = "/"
	httpCodeErr             = 500

	localhost         = "localhost:"
	websocketProtocol = "ws://"

	signedMsgPrefix = "vk"

	defaultWsPortOffset = 100 // The default offset between a Geth node's HTTP and websocket ports.
)

// ViewingKey is the packet of data sent to the enclave when storing a new viewing key.
type ViewingKey struct {
	publicKey *ecdsa.PublicKey
	signature []byte
}

// RunConfig contains the configuration required by StartWalletExtension.
type RunConfig struct {
	LocalNetwork      bool
	PrefundedAccounts []string
	StartPort         int
	UseFacade         bool
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

// TODO - Display error in browser if Metamask is not enabled (i.e. `ethereum` object is not available in-browser).
// TODO - Make node address configurable when not using local network.
// TODO - Add support for websockets on host server.

// StartWalletExtension starts the wallet extension and Obscuro facade, and optionally a local Ethereum network. It
// returns a handle to stop the wallet extension, Obscuro facade and local network nodes, if any were created.
func StartWalletExtension(config RunConfig) func() {
	nodeAddr := localhost + strconv.Itoa(config.StartPort+defaultWsPortOffset+2)

	var localNetwork *gethnetwork.GethNetwork
	if config.LocalNetwork {
		gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
		if err != nil {
			panic(err)
		}

		localNetwork = gethnetwork.NewGethNetwork(config.StartPort+2, config.StartPort+defaultWsPortOffset+2, gethBinaryPath, 1, 1, config.PrefundedAccounts)
		fmt.Println("Local Geth network started.")

		nodeAddr = localhost + strconv.Itoa(int(localNetwork.WebSocketPorts[0]))
	}

	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingKeyChannel := make(chan ViewingKey)

	// If we're using a facade, we point the wallet extension to the facade. Otherwise, we point it directly to the node.
	var walletExtensionForwardAddr string
	if config.UseFacade {
		walletExtensionForwardAddr = localhost + strconv.Itoa(config.StartPort+1)
	} else {
		walletExtensionForwardAddr = nodeAddr
	}

	walletExtensionAddr := localhost + strconv.Itoa(config.StartPort)
	walletExtension := NewWalletExtension(enclavePrivateKey, walletExtensionForwardAddr, viewingKeyChannel)

	var obscuroFacade *ObscuroFacade
	if config.UseFacade {
		obscuroFacade = NewObscuroFacade(enclavePrivateKey, websocketProtocol+nodeAddr, viewingKeyChannel)
		go obscuroFacade.Serve(walletExtensionForwardAddr)
		fmt.Println("Obscuro facade started.")
	}

	go walletExtension.Serve(walletExtensionAddr)
	fmt.Printf("Wallet extension started.\n💡 Visit %s/viewingkeys/ to generate an ephemeral viewing key. "+
		"Without a viewing key, you will not be able to decrypt the enclave's secure responses to your "+
		"eth_getBalance and eth_call requests.\n", walletExtensionAddr)

	// We return a handle to stop the components, including the local network nodes if any were created.
	shutdownFacadeAndExtension := func() {
		if obscuroFacade != nil {
			obscuroFacade.Shutdown()
		}
		walletExtension.Shutdown()
	}

	if !config.LocalNetwork {
		return shutdownFacadeAndExtension
	}
	return func() {
		localNetwork.StopNodes()
		shutdownFacadeAndExtension()
	}
}
