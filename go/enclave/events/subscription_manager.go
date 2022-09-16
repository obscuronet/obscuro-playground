package events

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/state"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/obscuronet/go-obscuro/go/common"
)

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	rpcEncryptionManager *rpc.EncryptionManager
	subscriptions        map[uuid.UUID]*common.LogSubscription
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager) *SubscriptionManager {
	return &SubscriptionManager{
		rpcEncryptionManager: rpcEncryptionManager,
		subscriptions:        map[uuid.UUID]*common.LogSubscription{},
	}
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
func (s *SubscriptionManager) AddSubscription(id uuid.UUID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	jsonSubscription, err := s.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	var subscriptions []common.LogSubscription
	if err := json.Unmarshal(jsonSubscription, &subscriptions); err != nil {
		return fmt.Errorf("could not unmarshall log subscription from JSON. Cause: %w", err)
	}

	if len(subscriptions) != 1 {
		return fmt.Errorf("expected a single log subscription, received %d", len(subscriptions))
	}
	subscription := subscriptions[0]

	err = s.rpcEncryptionManager.AuthenticateSubscriptionRequest(subscription)
	if err != nil {
		return err
	}

	s.subscriptions[id] = &subscription
	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
// TODO - #453 - Consider whether the deletion needs to be authenticated as well, to prevent attackers deleting subscriptions.
func (s *SubscriptionManager) RemoveSubscription(id uuid.UUID) {
	delete(s.subscriptions, id)
}

// FilterRelevantLogs filters out logs that are not subscribed too, and organises the logs by their subscribing ID.
func (s *SubscriptionManager) FilterRelevantLogs(logs []*types.Log, db *state.StateDB) map[uuid.UUID][]*types.Log {
	relevantLogs := map[uuid.UUID][]*types.Log{}

	for _, log := range logs {
		for subscriptionID, subscription := range s.subscriptions {
			logIsRelevant := isRelevant(log, subscription, db)
			if !logIsRelevant {
				continue
			}

			logsForSubID, found := relevantLogs[subscriptionID]
			if !found {
				relevantLogs[subscriptionID] = []*types.Log{log}
			} else {
				relevantLogs[subscriptionID] = append(logsForSubID, log)
			}
		}
	}

	return relevantLogs
}

// EncryptLogs encrypts each log with the appropriate viewing key.
// TODO - #453 - Encrypt logs, rather than just serialising them as JSON.
func (s *SubscriptionManager) EncryptLogs(logsBySubID map[uuid.UUID][]*types.Log) (map[uuid.UUID]common.EncryptedLogs, error) {
	result := map[uuid.UUID]common.EncryptedLogs{}
	for subID, logs := range logsBySubID {
		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			return nil, err
		}
		result[subID] = jsonLogs
	}
	return result, nil
}

// Indicates whether the log is relevant for the subscription. A lifecycle log is considered relevant to everyone.
func isRelevant(log *types.Log, sub *common.LogSubscription, db *state.StateDB) bool {
	var nonContractAddrs []*gethcommon.Address
	for _, topic := range log.Topics {
		addr := gethcommon.HexToAddress(topic.Hex())
		if db.GetCode(addr) == nil { // If there is code associated with the address, it's a contract address.
			nonContractAddrs = append(nonContractAddrs, &addr)
		}
	}

	if len(nonContractAddrs) == 0 {
		// All the topic addresses are contract addresses, so this is a log event, and is relevant to everyone.
		return true
	}

	for _, addr := range nonContractAddrs {
		if addr == sub.SubscriptionAccount.Account {
			return true
		}
	}
	return false
}
