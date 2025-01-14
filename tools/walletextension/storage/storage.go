package storage

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/sqlite"
)

// todo - pass the Context
type UserStorage interface {
	AddUser(userID []byte, privateKey []byte) error
	DeleteUser(userID []byte) error
	AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType viewingkey.SignatureType) error
	AddSessionKey(userID []byte, key common.GWSessionKey) error
	ActivateSessionKey(userID []byte, active bool) error
	RemoveSessionKey(userID []byte) error
	GetUser(userID []byte) (*common.GWUser, error)
	GetEncryptionKey() []byte
}

func New(dbType, dbConnectionURL, dbPath string, randomKey []byte, logger gethlog.Logger) (UserStorage, error) {
	var underlyingStorage UserStorage
	var err error
	switch dbType {
	case "sqlite":
		underlyingStorage, err = sqlite.NewSqliteDatabase(dbPath)
	case "cosmosDB":
		underlyingStorage, err = cosmosdb.NewCosmosDB(dbConnectionURL, randomKey)
	default:
		panic(fmt.Sprintf("unknown db type: %s", dbType))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize underlying storage: %w", err)
	}

	return NewUserStorageWithCache(underlyingStorage, logger)
}

// NewMetricsStorage is a factory function to create a MetricsStorage instance
func NewMetricsStorage(dbType, dbConnectionURL string) (*cosmosdb.MetricsStorageCosmosDB, error) {
	if dbType == "cosmosDB" {
		return cosmosdb.NewMetricsStorage(dbConnectionURL)
	}
	return nil, nil // Return nil for other database types
}
