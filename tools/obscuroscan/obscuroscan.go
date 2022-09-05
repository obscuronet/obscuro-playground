package obscuroscan

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/edgelesssys/ego/enclave"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/rpc"
)

const (
	pathAPI               = "/api"
	pathNumRollups        = "/numrollups/"
	pathNumTxs            = "/numtxs/"
	pathGetRollupTime     = "/rolluptime/"
	pathLatestRollups     = "/latestrollups/"
	pathLatestTxs         = "/latesttxs/"
	pathBlock             = "/block/"
	pathRollup            = "/rollup/"
	pathDecryptTxBlob     = "/decrypttxblob/"
	pathAttestation       = "/attestation/"
	pathAttestationReport = "/attestationreport/"
	pathRoot              = "/"

	staticDir   = "static"
	extDivider  = "."
	extHTML     = ".html"
	httpCodeErr = 500
)

//go:embed static
var staticFiles embed.FS

// Obscuroscan is a server that allows the monitoring of a running Obscuro network.
type Obscuroscan struct {
	server      *http.Server
	client      rpc.Client
	contractABI abi.ABI
}

// Identical to attestation.Report, but with the status mapped to a user-friendly string.
type attestationReportExternal struct {
	Data            []byte
	SecurityVersion uint
	Debug           bool
	UniqueID        []byte
	SignerID        []byte
	ProductID       []byte
	TCBStatus       string
}

func NewObscuroscan(address string) *Obscuroscan {
	client, err := rpc.NewNetworkClient(address)
	if err != nil {
		panic(err)
	}
	contractABI, err := abi.JSON(strings.NewReader(mgmtcontractlib.MgmtContractABI))
	if err != nil {
		panic("could not parse management contract ABI to decrypt rollups")
	}
	return &Obscuroscan{
		client:      client,
		contractABI: contractABI,
	}
}

// Serve listens for and serves Obscuroscan requests.
func (o *Obscuroscan) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathAPI+pathNumRollups, o.getNumRollups)            // Get the number of published rollups.
	serveMux.HandleFunc(pathAPI+pathNumTxs, o.getNumTransactions)           // Get the number of rolled-up transactions.
	serveMux.HandleFunc(pathAPI+pathGetRollupTime, o.getRollupTime)         // Get the average rollup time.
	serveMux.HandleFunc(pathAPI+pathLatestRollups, o.getLatestRollups)      // Get the latest rollup numbers.
	serveMux.HandleFunc(pathAPI+pathLatestTxs, o.getLatestTxs)              // Get the latest transaction hashes.
	serveMux.HandleFunc(pathAPI+pathRollup, o.getRollupByNumOrTxHash)       // Get the rollup given its number or the hash of a transaction it contains.
	serveMux.HandleFunc(pathAPI+pathBlock, o.getBlock)                      // Get the L1 block with the given number.
	serveMux.HandleFunc(pathAPI+pathDecryptTxBlob, o.decryptTxBlob)         // Decrypt a transaction blob.
	serveMux.HandleFunc(pathAPI+pathAttestation, o.attestation)             // Retrieve the node's attestation.
	serveMux.HandleFunc(pathAPI+pathAttestationReport, o.attestationReport) // Retrieve the node's attestation report.

	// Serves the web assets for the user interface.
	staticFileFS, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	staticFileFilesystem := http.FileServer(http.FS(staticFileFS))
	serveMux.Handle(pathRoot, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If we get a request without an extension (other than for the root), we tack on ".html".
		if r.URL.Path != pathRoot && !strings.Contains(r.URL.Path, extDivider) {
			r.URL.Path += extHTML
		}
		staticFileFilesystem.ServeHTTP(w, r)
	}))

	o.server = &http.Server{Addr: hostAndPort, Handler: serveMux, ReadHeaderTimeout: 10 * time.Second}
	err = o.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (o *Obscuroscan) Shutdown() {
	if o.server != nil {
		err := o.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down Obscuroscan. Cause: %s", err)
		}
	}
}

// Retrieves the number of published rollups.
func (o *Obscuroscan) getNumRollups(resp http.ResponseWriter, _ *http.Request) {
	numOfRollups, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}

	numOfRollupsStr := strconv.Itoa(int(numOfRollups))
	_, err = resp.Write([]byte(numOfRollupsStr))
	if err != nil {
		log.Error("could not return number of rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch number of rollups.")
		return
	}
}

// Retrieves the total number of transactions.
func (o *Obscuroscan) getNumTransactions(resp http.ResponseWriter, _ *http.Request) {
	var numTransactions *big.Int
	err := o.client.Call(&numTransactions, rpc.RPCGetTotalTxs)
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch total transactions.")
		return
	}

	_, err = resp.Write([]byte(numTransactions.String()))
	if err != nil {
		log.Error("could not return total number of transactions to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch total transactions.")
		return
	}
}

// Retrieves the average rollup time, as (time last rollup - time first rollup)/number of rollups
func (o *Obscuroscan) getRollupTime(resp http.ResponseWriter, _ *http.Request) {
	numLatestRollup, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	firstRollup, err := o.getRollupByNumber(0)
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	latestRollup, err := o.getRollupByNumber(int(numLatestRollup))
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}

	avgRollupTime := float64(latestRollup.Header.Time-firstRollup.Header.Time) / float64(numLatestRollup)
	_, err = resp.Write([]byte(fmt.Sprintf("%.2f", avgRollupTime)))
	if err != nil {
		log.Error("could not return average rollup time to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch average rollup time.")
		return
	}
}

// Retrieves the last five rollup numbers.
func (o *Obscuroscan) getLatestRollups(resp http.ResponseWriter, _ *http.Request) {
	latestRollupNum, err := o.getLatestRollupNumber()
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}

	// We walk the chain of rollups, getting the number for the most recent five.
	rollupNums := make([]string, 5)
	for idx := 0; idx < 5; idx++ {
		rollupNum := int(latestRollupNum) - idx
		if rollupNum < 0 {
			// If there are less than five rollups, we return an N/A.
			rollupNums[idx] = "N/A"
		} else {
			rollupNums[idx] = strconv.Itoa(rollupNum)
		}
	}

	jsonRollupNums, err := json.Marshal(rollupNums)
	if err != nil {
		log.Error("could not return latest rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
	_, err = resp.Write(jsonRollupNums)
	if err != nil {
		log.Error("could not return latest rollups to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest rollups.")
		return
	}
}

// Retrieves the last five transaction hashes.
func (o *Obscuroscan) getLatestTxs(resp http.ResponseWriter, _ *http.Request) {
	numTransactions := 5

	var txHashes []gethcommon.Hash
	err := o.client.Call(&txHashes, rpc.RPCGetLatestTxs, numTransactions)
	if err != nil {
		log.Error(err.Error())
		logAndSendErr(resp, "Could not fetch latest transactions.")
	}

	// We convert the hashes to strings and pad with N/As as needed.
	txHashStrings := make([]string, numTransactions)
	for idx := 0; idx < numTransactions; idx++ {
		if idx < len(txHashes) {
			txHashStrings[idx] = txHashes[idx].String()
		} else {
			txHashStrings[idx] = "N/A"
		}
	}

	jsonTxHashes, err := json.Marshal(txHashStrings)
	if err != nil {
		log.Error("could not return latest transaction hashes to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}
	_, err = resp.Write(jsonTxHashes)
	if err != nil {
		log.Error("could not return latest transaction hashes to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch latest transactions.")
		return
	}
}

// Retrieves the L1 block header with the given number.
func (o *Obscuroscan) getBlock(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	blockHashStr := buffer.String()
	blockHash := gethcommon.HexToHash(blockHashStr)

	var blockHeader *types.Header
	err = o.client.Call(&blockHeader, rpc.RPCGetBlockHeaderByHash, blockHash)
	if err != nil {
		log.Error("could not retrieve block with hash %s. Cause: %s", blockHash, err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}

	jsonBlock, err := json.Marshal(blockHeader)
	if err != nil {
		log.Error("could not return block to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
	_, err = resp.Write(jsonBlock)
	if err != nil {
		log.Error("could not return block to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch block.")
		return
	}
}

// Retrieves a rollup given its number or the hash of a transaction it contains.
func (o *Obscuroscan) getRollupByNumOrTxHash(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}

	var rollup *common.ExtRollup
	if strings.HasPrefix(buffer.String(), "0x") {
		// A "0x" prefix indicates that we should retrieve the rollup by transaction hash.
		txHash := gethcommon.HexToHash(buffer.String())

		err = o.client.Call(&rollup, rpc.RPCGetRollupForTx, txHash)
		if err != nil {
			log.Error("could not retrieve rollup. Cause: %s", err)
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
	} else {
		// Otherwise, we treat the input as a rollup number.
		rollupNumber, err := strconv.Atoi(buffer.String())
		if err != nil {
			log.Error("could not parse \"%s\" as an integer", buffer.String())
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
		rollup, err = o.getRollupByNumber(rollupNumber)
		if err != nil {
			log.Error(err.Error())
			logAndSendErr(resp, "Could not fetch rollup.")
			return
		}
	}

	jsonRollup, err := json.Marshal(rollup)
	if err != nil {
		log.Error("could not return rollup to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}
	_, err = resp.Write(jsonRollup)
	if err != nil {
		log.Error("could not return rollup to client. Cause: %s", err)
		logAndSendErr(resp, "Could not fetch rollup.")
		return
	}
}

// Decrypts the provided transaction blob using the provided key.
// TODO - Use the passed-in key, rather than a hardcoded enclave key.
func (o *Obscuroscan) decryptTxBlob(resp http.ResponseWriter, req *http.Request) {
	body := req.Body
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		log.Error("could not read request body. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	jsonTxs, err := decryptTxBlob(buffer.Bytes())
	if err != nil {
		log.Error("could not decrypt transaction blob. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}

	_, err = resp.Write(jsonTxs)
	if err != nil {
		log.Error("could not write decrypted transactions to client. Cause: %s", err)
		logAndSendErr(resp, "Could not decrypt transaction blob.")
		return
	}
}

// Retrieves the node's attestation.
func (o *Obscuroscan) attestation(resp http.ResponseWriter, _ *http.Request) {
	var attestation *common.AttestationReport
	err := o.client.Call(&attestation, rpc.RPCAttestation)
	if err != nil {
		log.Error("could not retrieve node's attestation. Cause: %s", err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}

	jsonAttestation, err := json.Marshal(attestation)
	if err != nil {
		log.Error("could not convert node's attestation to JSON. Cause: %s", err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}
	_, err = resp.Write(jsonAttestation)
	if err != nil {
		log.Error("could not return JSON attestation to client. Cause: %s", err)
		logAndSendErr(resp, "Could not retrieve node's attestation.")
		return
	}
}

// Retrieves the node's attestation report.
func (o *Obscuroscan) attestationReport(resp http.ResponseWriter, _ *http.Request) {
	var attestation *common.AttestationReport
	err := o.client.Call(&attestation, rpc.RPCAttestation)
	if err != nil {
		log.Error("could not retrieve node's attestation. Cause: %s", err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}

	// If DCAP isn't set up, verifying the report will send a SIGSYS signal. We catch this, so that it doesn't crash the program.
	sigChannel := make(chan os.Signal, 1)
	defer signal.Stop(sigChannel)
	signal.Notify(sigChannel, syscall.SIGSYS)
	attestationReport, err := enclave.VerifyRemoteReport(attestation.Report)
	signal.Stop(sigChannel)

	if err != nil {
		log.Error("could not verify node's attestation. Cause: %s", err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}

	attestationReportExt := attestationReportExternal{
		Data:            attestationReport.Data,
		SecurityVersion: attestationReport.SecurityVersion,
		Debug:           attestationReport.Debug,
		UniqueID:        attestationReport.UniqueID,
		SignerID:        attestationReport.SignerID,
		ProductID:       attestationReport.ProductID,
		TCBStatus:       attestationReport.TCBStatus.String(),
	}

	jsonAttestationReport, err := json.Marshal(attestationReportExt)
	if err != nil {
		log.Error("could not convert node's attestation report to JSON. Cause: %s", err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}
	_, err = resp.Write(jsonAttestationReport)
	if err != nil {
		log.Error("could not return JSON attestation report to client. Cause: %s", err)
		logAndSendErr(resp, "Could not verify node's attestation.")
		return
	}
}

// Returns the number of the latest rollup.
func (o *Obscuroscan) getLatestRollupNumber() (int64, error) {
	var rollupHeader *common.Header
	err := o.client.Call(&rollupHeader, rpc.RPCGetCurrentRollupHead)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve head rollup. Cause: %w", err)
	}

	latestRollupNum := rollupHeader.Number.Int64()
	return latestRollupNum, nil
}

// Parses numberStr as a number and returns the rollup with that number.
func (o *Obscuroscan) getRollupByNumber(rollupNumber int) (*common.ExtRollup, error) {
	// TODO - If required, consolidate the two calls below into a single RPCGetRollupByNumber call to minimise round trips.
	var rollupHeader *common.Header
	err := o.client.Call(&rollupHeader, rpc.RPCGetRollupHeaderByNumber, rollupNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with number %d. Cause: %w", rollupNumber, err)
	}

	rollupHash := rollupHeader.Hash()
	if rollupHash == (gethcommon.Hash{}) {
		return nil, fmt.Errorf("rollup was retrieved but hash was nil")
	}

	var rollup *common.ExtRollup
	err = o.client.Call(&rollup, rpc.RPCGetRollup, rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup. Cause: %w", err)
	}
	return rollup, nil
}

// Decrypts the transaction blob and returns it as JSON.
func decryptTxBlob(encryptedTxBytesBase64 []byte) ([]byte, error) {
	encryptedTxBytes, err := base64.StdEncoding.DecodeString(string(encryptedTxBytesBase64))
	if err != nil {
		return nil, fmt.Errorf("could not decode encrypted transaction blob from Base64. Cause: %w", err)
	}

	key := gethcommon.Hex2Bytes(crypto.RollupEncryptionKeyHex)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not initialise AES cipher for enclave rollup key. Cause: %w", err)
	}
	transactionCipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not initialise wrapper for AES cipher for enclave rollup key. Cause: %w", err)
	}

	// The nonce is prepended to the ciphertext.
	nonce := encryptedTxBytes[0:crypto.NonceLength]
	ciphertext := encryptedTxBytes[crypto.NonceLength:]
	encodedTxs, err := transactionCipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt encrypted L2 transactions. Cause: %w", err)
	}

	var cleartextTxs []*common.L2Tx
	if err = rlp.DecodeBytes(encodedTxs, &cleartextTxs); err != nil {
		return nil, fmt.Errorf("could not decode encoded L2 transactions. Cause: %w", err)
	}

	jsonRollup, err := json.Marshal(cleartextTxs)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt transaction blob. Cause: %w", err)
	}

	return jsonRollup, nil
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}
