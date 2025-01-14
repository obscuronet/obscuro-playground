package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/storage/database/cosmosdb"
)

const (
	// Persistence intervals (how often metrics are saved to CosmosDB)
	MetricsPersistInterval = 10 * time.Minute

	// Cleanup intervals (how often inactive users are cleaned up)
	InactiveUserCleanupInterval = 1 * time.Hour

	// Activity thresholds
	UserInactivityThreshold = 30 * 24 * time.Hour // 30 days
	MonthlyActiveUserWindow = 30 * 24 * time.Hour // 30 days
)

// Metrics interface defines the metrics operations
type Metrics interface {
	RecordNewUser()
	RecordAccountRegistered()
	RecordUserActivity(anonymousID string)
	GetTotalUsers() uint64
	GetTotalAccountsRegistered() uint64
	GetMonthlyActiveUsers() int
	Stop()
}

type MetricsTracker struct {
	totalUsers         atomic.Uint64
	accountsRegistered atomic.Uint64
	activeUsers        map[string]time.Time // key is double-hashed userID
	activeUserLock     sync.RWMutex
	storage            *cosmosdb.MetricsStorageCosmosDB
	persistTicker      *time.Ticker
}

func NewMetricsTracker(storage *cosmosdb.MetricsStorageCosmosDB) Metrics {
	mt := &MetricsTracker{
		activeUsers:   make(map[string]time.Time),
		storage:       storage,
		persistTicker: time.NewTicker(MetricsPersistInterval),
	}

	// Load existing metrics
	if metrics, err := storage.LoadMetrics(); err == nil {
		mt.totalUsers.Store(metrics.TotalUsers)
		mt.accountsRegistered.Store(metrics.AccountsRegistered)

		mt.activeUserLock.Lock()
		for hashedUserID, timestamp := range metrics.ActiveUsers {
			if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
				mt.activeUsers[hashedUserID] = t
			}
		}
		mt.activeUserLock.Unlock()
	}

	// Start cleanup routine for inactive users
	go mt.cleanupInactiveUsers()
	go mt.persistMetrics()

	return mt
}

// hashUserID creates a double-hashed version of the userID
func (mt *MetricsTracker) hashUserID(userID []byte) string {
	// First hash
	firstHash := sha256.Sum256(userID)
	// Second hash
	secondHash := sha256.Sum256(firstHash[:])
	return hex.EncodeToString(secondHash[:])
}

func (mt *MetricsTracker) RecordNewUser() {
	mt.totalUsers.Add(1)
}

// RecordAccountRegistered increments the total number of registered accounts
func (mt *MetricsTracker) RecordAccountRegistered() {
	mt.accountsRegistered.Add(1)
}

// RecordUserActivity updates the last activity timestamp for a user
func (mt *MetricsTracker) RecordUserActivity(anonymousID string) {
	hashedUserID := mt.hashUserID([]byte(anonymousID))

	mt.activeUserLock.Lock()
	mt.activeUsers[hashedUserID] = time.Now()
	mt.activeUserLock.Unlock()
}

// GetTotalUsers returns the total number of registered users
func (mt *MetricsTracker) GetTotalUsers() uint64 {
	return mt.totalUsers.Load()
}

// GetTotalAccountsRegistered returns the total number of registered accounts
func (mt *MetricsTracker) GetTotalAccountsRegistered() uint64 {
	return mt.accountsRegistered.Load()
}

// GetMonthlyActiveUsers returns the number of users active in the last 30 days
func (mt *MetricsTracker) GetMonthlyActiveUsers() int {
	mt.activeUserLock.RLock()
	defer mt.activeUserLock.RUnlock()

	count := 0
	activeThreshold := time.Now().Add(-MonthlyActiveUserWindow)

	for _, lastActive := range mt.activeUsers {
		if lastActive.After(activeThreshold) {
			count++
		}
	}
	return count
}

// persistMetrics periodically saves metrics to CosmosDB
func (mt *MetricsTracker) persistMetrics() {
	for range mt.persistTicker.C {
		mt.saveMetrics()
	}
}

func (mt *MetricsTracker) saveMetrics() {
	mt.activeUserLock.RLock()
	activeUsersMap := make(map[string]string)
	for hashedUserID, timestamp := range mt.activeUsers {
		activeUsersMap[hashedUserID] = timestamp.UTC().Format(time.RFC3339)
	}
	mt.activeUserLock.RUnlock()

	metrics := &cosmosdb.MetricsDocument{
		ID:                 cosmosdb.METRICS_DOC_ID,
		TotalUsers:         mt.totalUsers.Load(),
		AccountsRegistered: mt.accountsRegistered.Load(),
		ActiveUsers:        activeUsersMap,
	}

	if err := mt.storage.SaveMetrics(metrics); err != nil {
		// Either log the error properly or return it
		log.Printf("Failed to persist metrics: %v", err)
	}
}

func (mt *MetricsTracker) cleanupInactiveUsers() {
	ticker := time.NewTicker(InactiveUserCleanupInterval)
	for range ticker.C {
		mt.activeUserLock.Lock()
		inactiveThreshold := time.Now().Add(-UserInactivityThreshold)

		for userID, lastActive := range mt.activeUsers {
			if lastActive.Before(inactiveThreshold) {
				delete(mt.activeUsers, userID)
			}
		}
		mt.activeUserLock.Unlock()
	}
}

// Stop cleanly stops the metrics tracker
func (mt *MetricsTracker) Stop() {
	mt.persistTicker.Stop()
	mt.saveMetrics() // Final save before stopping
}

// NoOpMetricsTracker implements Metrics interface but does nothing
type NoOpMetricsTracker struct{}

func NewNoOpMetricsTracker() Metrics {
	return &NoOpMetricsTracker{}
}

func (mt *NoOpMetricsTracker) RecordNewUser()                     {}
func (mt *NoOpMetricsTracker) RecordAccountRegistered()           {}
func (mt *NoOpMetricsTracker) RecordUserActivity(string)          {}
func (mt *NoOpMetricsTracker) GetTotalUsers() uint64              { return 0 }
func (mt *NoOpMetricsTracker) GetTotalAccountsRegistered() uint64 { return 0 }
func (mt *NoOpMetricsTracker) GetMonthlyActiveUsers() int         { return 0 }
func (mt *NoOpMetricsTracker) Stop()                              {}
