package lib

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/valyala/fasthttp"
)

type OGLib struct {
	httpURL string
	wsURL   string
	userID  []byte
}

func NewObscuroGatewayLibrary(httpURL, wsURL string) *OGLib {
	return &OGLib{
		httpURL: httpURL,
		wsURL:   wsURL,
	}
}

func (o *OGLib) UserID() string {
	return string(o.userID)
}

func (o *OGLib) Join() error {
	// todo move this to stdlib
	statusCode, userID, err := fasthttp.Get(nil, fmt.Sprintf("%s/v1/join/", o.httpURL))
	if err != nil || statusCode != 200 {
		return fmt.Errorf(fmt.Sprintf("Failed to get userID. Status code: %d, err: %s", statusCode, err))
	}
	o.userID = userID
	return nil
}

func (o *OGLib) RegisterAccount(pk *ecdsa.PrivateKey, addr gethcommon.Address) error {
	// create the registration message
	rawMessage, err := viewingkey.GenerateAuthenticationEIP712RawData(string(o.userID))
	if err != nil {
		return err
	}

	messageHash := crypto.Keccak256(rawMessage)
	sig, err := crypto.Sign(messageHash, pk)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	sig[64] += 27
	signature := "0x" + hex.EncodeToString(sig)
	payload := fmt.Sprintf("{\"signature\": \"%s\", \"address\": \"%s\"}", signature, addr.Hex())

	// issue the registration message
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		o.httpURL+"/v1/authenticate/?u="+string(o.userID),
		strings.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("unable to create request - %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to issue request - %w", err)
	}

	defer response.Body.Close()
	r, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read response - %w", err)
	}
	if string(r) != "success" {
		return fmt.Errorf("expected success, got %s", string(r))
	}
	return nil
}

func (o *OGLib) HTTP() string {
	return fmt.Sprintf("%s/v1/?u=%s", o.httpURL, o.userID)
}

func (o *OGLib) WS() string {
	return fmt.Sprintf("%s/v1/?u=%s", o.wsURL, o.userID)
}
