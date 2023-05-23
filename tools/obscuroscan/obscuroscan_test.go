package obscuroscan

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"
)

func TestCanDecryptTxBlob(t *testing.T) {
	txs := []*common.L2Tx{datagenerator.CreateL2Tx(), datagenerator.CreateL2Tx()}

	txsJSONBytes, err := decryptTxBlob(generateEncryptedTxBlob(txs))
	if err != nil {
		t.Fatalf("transaction blob decryption failed. Cause: %s", err)
	}

	expectedTxsJSONBytes, err := json.Marshal(txs)
	if err != nil {
		t.Fatalf("marshalling transactions to JSON failed. Cause: %s", err)
	}

	if string(expectedTxsJSONBytes) != string(txsJSONBytes) {
		t.Fatalf("expected %s, got %s", string(expectedTxsJSONBytes), string(txsJSONBytes))
	}
}

func TestThrowsIfEncryptedRollupIsInvalid(t *testing.T) {
	_, err := decryptTxBlob([]byte("invalid_tx_blob"))
	if err == nil {
		t.Fatal("did not error on invalid transaction blob")
	}
}

// Generates an encrypted transaction blob in Base64 encoding.
func generateEncryptedTxBlob(txs []*common.L2Tx) []byte {
	rollup := core.Batch{Header: &common.BatchHeader{}, Transactions: txs}
	extB, err := rollup.ToExtBatch(crypto.NewDataEncryptionService(nil), crypto.NewBrotliDataCompressionService())
	if err != nil {
		panic(err)
	}
	txBlob := extB.EncryptedTxBlob
	return []byte(base64.StdEncoding.EncodeToString(txBlob))
}

func TestObscuroscan_getRollupByNumOrTxHash(t *testing.T) {
	logger := gethlog.Logger.New(gethlog.Root())
	ob := NewObscuroscan("http://testnet.obscuroscan.io", logger)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.Method = http.MethodOptions
	resp := httptest.NewRecorder()
	ob.getRollupByNumOrTxHash(resp, req)
	if resp.Header().Get(httputil.CorsAllowOrigin) != httputil.OriginAll {
		t.Fatal("CORS Allow Origin not set.")
	}
	if resp.Header().Get(httputil.CorsAllowMethods) != httputil.ReqOptions {
		t.Fatal("CORS Allow Methods not set.")
	}
	if resp.Header().Get(httputil.CorsAllowHeaders) != httputil.CorsHeaders {
		t.Fatal("CORS Allow Headers not set.")
	}
}
