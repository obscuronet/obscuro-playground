package events

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/eth/filters"

	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/ethereum/go-ethereum/core/state"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// The leading zero bytes in a hash indicating that it is possibly an address, since it only has 20 bytes of data.
	zeroBytesHex = "000000000000000000000000"
)

// TODO - Ensure chain reorgs are handled gracefully.

// SubscriptionManager manages the creation/deletion of subscriptions, and the filtering and encryption of logs for
// active subscriptions.
type SubscriptionManager struct {
	rpcEncryptionManager *rpc.EncryptionManager
	storage              db.Storage

	subscriptions     map[gethrpc.ID]*common.LogSubscription
	subscriptionMutex *sync.RWMutex
}

func NewSubscriptionManager(rpcEncryptionManager *rpc.EncryptionManager, storage db.Storage) *SubscriptionManager {
	return &SubscriptionManager{
		rpcEncryptionManager: rpcEncryptionManager,
		storage:              storage,

		subscriptions:     map[gethrpc.ID]*common.LogSubscription{},
		subscriptionMutex: &sync.RWMutex{},
	}
}

// AddSubscription adds a log subscription to the enclave under the given ID, provided the request is authenticated
// correctly. If there is an existing subscription with the given ID, it is overwritten.
func (s *SubscriptionManager) AddSubscription(id gethrpc.ID, encryptedSubscription common.EncryptedParamsLogSubscription) error {
	encodedSubscription, err := s.rpcEncryptionManager.DecryptBytes(encryptedSubscription)
	if err != nil {
		return fmt.Errorf("could not decrypt params in eth_subscribe logs request. Cause: %w", err)
	}

	var subscription common.LogSubscription
	if err = rlp.DecodeBytes(encodedSubscription, &subscription); err != nil {
		return fmt.Errorf("could not decocde log subscription from RLP. Cause: %w", err)
	}

	err = s.rpcEncryptionManager.AuthenticateSubscriptionRequest(subscription)
	if err != nil {
		return err
	}

	// For subscriptions, only the Topics and Addresses fields of the filter are applied.
	subscription.Filter.BlockHash = nil
	subscription.Filter.ToBlock = nil
	// We set this to the current rollup height, so that historical logs aren't returned.
	subscription.Filter.FromBlock = big.NewInt(0).Add(s.storage.FetchHeadRollup().Number(), big.NewInt(1))

	s.subscriptionMutex.Lock()
	s.subscriptions[id] = &subscription
	s.subscriptionMutex.Unlock()
	return nil
}

// RemoveSubscription removes the log subscription with the given ID from the enclave. If there is no subscription with
// the given ID, nothing is deleted.
func (s *SubscriptionManager) RemoveSubscription(id gethrpc.ID) {
	s.subscriptionMutex.Lock()
	delete(s.subscriptions, id)
	s.subscriptionMutex.Unlock()
}

// FilteredLogsHead retrieves the logs for the head of the chain and filters out irrelevant logs.
func (s *SubscriptionManager) FilteredLogsHead(account *gethcommon.Address, filter *filters.FilterCriteria) ([]*types.Log, error) {
	headBlockHash := s.storage.FetchHeadBlock().Hash()
	_, logs, found := s.storage.FetchBlockState(headBlockHash)
	if !found {
		log.Error("could not retrieve logs for head state. Something is wrong")
		return nil, fmt.Errorf("could not retrieve logs for head state. Something is wrong")
	}
	return s.FilteredLogs(logs, s.storage.FetchHeadRollup().Hash(), account, filter), nil
}

// FilteredLogs filters out irrelevant logs from the list provided.
func (s *SubscriptionManager) FilteredLogs(logs []*types.Log, rollupHash common.L2RootHash, account *gethcommon.Address, filter *filters.FilterCriteria) []*types.Log {
	filteredLogs := []*types.Log{}
	stateDB := s.storage.CreateStateDB(rollupHash)

	for _, logItem := range logs {
		userAddrs := getUserAddrs(logItem, stateDB)
		if isRelevant(userAddrs, account) && !isFilteredOut(logItem, filter) {
			filteredLogs = append(filteredLogs, logItem)
		}
	}

	return filteredLogs
}

// FilteredSubscribedLogs filters out irrelevant logs, those that are not subscribed to, and those a subscription has
// seen before, and organises them by their subscribing ID.
func (s *SubscriptionManager) FilteredSubscribedLogs(logs []*types.Log, rollupHash common.L2RootHash) map[gethrpc.ID][]*types.Log {
	relevantLogs := map[gethrpc.ID][]*types.Log{}

	// If there are no subscriptions, we do not need to do any processing.
	if len(s.subscriptions) == 0 {
		return relevantLogs
	}

	stateDB := s.storage.CreateStateDB(rollupHash)

	lastSeenRollupByID := map[gethrpc.ID]uint64{}

	for _, logItem := range logs {
		userAddrs := getUserAddrs(logItem, stateDB)
		rollupNumber := logItem.BlockNumber

		// We check whether the log is relevant to each subscription.
		s.subscriptionMutex.RLock()
		for subscriptionID, subscription := range s.subscriptions {
			if isRelevant(userAddrs, subscription.Account) && !isFilteredOut(logItem, subscription.Filter) && rollupNumber > subscription.LastSeenRollup {
				if relevantLogs[subscriptionID] == nil {
					relevantLogs[subscriptionID] = []*types.Log{}
				}

				relevantLogs[subscriptionID] = append(relevantLogs[subscriptionID], logItem)

				// We update the last rollup number for the subscription if this is the highest one we've seen so far.
				currentLatestSeenRollup, exists := lastSeenRollupByID[subscriptionID]
				if !exists || rollupNumber > currentLatestSeenRollup {
					lastSeenRollupByID[subscriptionID] = rollupNumber
				}
			}
		}
		s.subscriptionMutex.RUnlock()
	}

	// We update the last seen rollup for each subscription. We must do this in a separate loop, as if we update it
	// as we go, we may miss a set of logs if we process the logs for rollup N before we process those for rollup N-1.
	s.subscriptionMutex.RLock()
	for subscriptionID, subscription := range s.subscriptions {
		newLatestSeenRollup, exists := lastSeenRollupByID[subscriptionID]
		if exists && newLatestSeenRollup > subscription.LastSeenRollup {
			subscription.LastSeenRollup = newLatestSeenRollup
		}
	}
	s.subscriptionMutex.RUnlock()

	return relevantLogs
}

// EncryptLogs encrypts each log with the appropriate viewing key.
func (s *SubscriptionManager) EncryptLogs(logsByID map[gethrpc.ID][]*types.Log) (map[gethrpc.ID][]byte, error) {
	encryptedLogsByID := map[gethrpc.ID][]byte{}

	for subID, logs := range logsByID {
		subscription, found := s.subscriptions[subID]
		if !found {
			continue // The subscription has been removed, so there's no need to return anything.
		}

		jsonLogs, err := json.Marshal(logs)
		if err != nil {
			return nil, fmt.Errorf("could not marshal logs to JSON. Cause: %w", err)
		}

		encryptedLogs, err := s.rpcEncryptionManager.EncryptWithViewingKey(*subscription.Account, jsonLogs)
		if err != nil {
			return nil, err
		}

		encryptedLogsByID[subID] = encryptedLogs
	}

	return encryptedLogsByID, nil
}

// Extracts the (potential) user addresses from the topics. If there is no code associated with an address, it's a user
// address.
func getUserAddrs(log *types.Log, db *state.StateDB) []string {
	var userAddrs []string //nolint:prealloc

	for _, topic := range log.Topics {
		// Since addresses are 20 bytes long, while hashes are 32, only topics with 12 leading zero bytes can
		// (potentially) be user addresses.
		topicHex := topic.Hex()
		if topicHex[2:len(zeroBytesHex)+2] != zeroBytesHex {
			continue
		}

		// If there is code associated with an address, it is not a user address.
		addr := gethcommon.HexToAddress(topicHex)
		if db.GetCode(addr) != nil {
			continue
		}

		userAddrs = append(userAddrs, addr.Hex())
	}

	return userAddrs
}

// Indicates whether the log is relevant for the subscription.
func isRelevant(userAddrs []string, account *gethcommon.Address) bool {
	// If there are no potential user addresses, this is a lifecycle event, and is therefore relevant to everyone.
	if len(userAddrs) == 0 {
		return true
	}

	for _, addr := range userAddrs {
		if addr == account.Hex() {
			return true
		}
	}

	return false
}

// Applies `filterLogs`, below, to determine whether the log should be filtered out based on the user's subscription criteria.
func isFilteredOut(log *types.Log, filterCriteria *filters.FilterCriteria) bool {
	filteredLogs := filterLogs([]*types.Log{log}, filterCriteria.FromBlock, filterCriteria.ToBlock, filterCriteria.Addresses, filterCriteria.Topics)
	return len(filteredLogs) == 0
}

// Lifted from eth/filters/filter.go in the go-ethereum repository.
// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*types.Log, fromBlock, toBlock *big.Int, addresses []gethcommon.Address, topics [][]gethcommon.Hash) []*types.Log { //nolint:gocognit
	var ret []*types.Log
Logs:
	for _, logItem := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > logItem.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < logItem.BlockNumber {
			continue
		}

		if len(addresses) > 0 && !includes(addresses, logItem.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(logItem.Topics) {
			continue
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if logItem.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, logItem)
	}
	return ret
}

// Lifted from eth/filters/filter.go in the go-ethereum repository.
func includes(addresses []gethcommon.Address, a gethcommon.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}
