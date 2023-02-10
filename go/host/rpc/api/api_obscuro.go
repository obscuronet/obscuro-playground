package api

import (
	"github.com/obscuronet/go-obscuro/go/common/host"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host host.Host
}

func NewObscuroAPI(host host.Host) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// AddViewingKey stores the viewing key on the enclave.
func (api *ObscuroAPI) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	return api.host.EnclaveClient().AddViewingKey(viewingKeyBytes, signature)
}

// Health returns the health status of obscuro host + enclave + db
func (api *ObscuroAPI) Health() (*host.HealthCheck, error) {
	return api.host.HealthCheck()
}
