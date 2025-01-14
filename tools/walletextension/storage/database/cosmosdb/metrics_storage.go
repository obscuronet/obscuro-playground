package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	METRICS_CONTAINER_NAME = "metrics"
	METRICS_DOC_ID         = "global_metrics"
)

type MetricsDocument struct {
	ID                 string            `json:"id"`
	TotalUsers         uint64            `json:"totalUsers"`
	AccountsRegistered uint64            `json:"accountsRegistered"`
	ActiveUsers        map[string]string `json:"activeUsers"` // double-hashed userID -> ISO timestamp
	ActiveUsersCount   int               `json:"activeUsersCount"`
	LastUpdated        string            `json:"lastUpdated"`
}

// MetricsStorageCosmosDB handles metrics persistence in CosmosDB
type MetricsStorageCosmosDB struct {
	client           *azcosmos.Client
	metricsContainer *azcosmos.ContainerClient
}

func NewMetricsStorage(connectionString string) (*MetricsStorageCosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	ctx := context.Background()

	// Ensure database exists
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	metricsContainer, err := client.NewContainer(DATABASE_NAME, METRICS_CONTAINER_NAME)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics container: %w", err)
	}

	return &MetricsStorageCosmosDB{
		client:           client,
		metricsContainer: metricsContainer,
	}, nil
}

func (m *MetricsStorageCosmosDB) LoadMetrics() (*MetricsDocument, error) {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(METRICS_DOC_ID)

	response, err := m.metricsContainer.ReadItem(ctx, partitionKey, METRICS_DOC_ID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// Initialize with empty metrics if not found
			return &MetricsDocument{
				ID:          METRICS_DOC_ID,
				ActiveUsers: make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var doc MetricsDocument
	if err := json.Unmarshal(response.Value, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (m *MetricsStorageCosmosDB) SaveMetrics(metrics *MetricsDocument) error {
	ctx := context.Background()
	partitionKey := azcosmos.NewPartitionKeyString(METRICS_DOC_ID)

	// Calculate active users count and clean up inactive users
	activeThreshold := time.Now().Add(-30 * 24 * time.Hour) // 30 days
	activeCount := 0
	activeUsersMap := make(map[string]string)

	for userID, timestampStr := range metrics.ActiveUsers {
		if timestamp, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			if timestamp.After(activeThreshold) {
				activeCount++
				activeUsersMap[userID] = timestampStr
			}
		}
	}

	metrics.ActiveUsers = activeUsersMap // Only keep active users
	metrics.ActiveUsersCount = activeCount
	metrics.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	docJSON, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = m.metricsContainer.UpsertItem(ctx, partitionKey, docJSON, nil)
	return err
}

// NoOpMetricsStorage is a no-op implementation of metrics storage
type noOpMetricsStorage struct{}

// MetricsStorage interface defines the metrics storage operations
type MetricsStorage interface {
	LoadMetrics() (*MetricsDocument, error)
	SaveMetrics(*MetricsDocument) error
}

func NewNoOpMetricsStorage() MetricsStorage {
	return &noOpMetricsStorage{}
}

func (n *noOpMetricsStorage) LoadMetrics() (*MetricsDocument, error) {
	return &MetricsDocument{
		ActiveUsers: make(map[string]string),
	}, nil
}

func (n *noOpMetricsStorage) SaveMetrics(*MetricsDocument) error {
	return nil
}
