package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
)

// CreateTLSHTTPClient provides a basic http client prepared with a trusted CA cert
func CreateTLSHTTPClient(caCertPEM string) (*http.Client, error) {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(caCertPEM)); !ok {
		return nil, fmt.Errorf("failed to append to CA cert from caCertPEM=%s", caCertPEM)
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}, nil
}

// ExecuteHTTPReq executes an HTTP request:
// * returns an error if request fails or if the response code was outside the range 200-299
// * returns response body as bytes if there was a response body
func ExecuteHTTPReq(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed url=%s - %w", req.URL.String(), err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg []byte
		_, err := resp.Body.Read(msg)
		if err != nil {
			return nil, fmt.Errorf("req failed url=%s, statusCode=%d, failed to read status text", req.URL.String(), resp.StatusCode)
		}
		return nil, fmt.Errorf("req failed url=%s status: %d %s", req.URL.String(), resp.StatusCode, msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// success status code but no body, ignoring error and returning no bytes
		return []byte{}, nil //nolint:nilerr
	}
	return body, nil
}
