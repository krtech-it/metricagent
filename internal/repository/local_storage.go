package repository

import (
	"context"
	"database/sql"
	"fmt"
	models "github.com/krtech-it/metricagent/internal/model"
)

type Storage interface {
	Update(ctx context.Context, metric *models.Metrics) error
	Create(ctx context.Context, metric *models.Metrics) error
	Get(ctx context.Context, id string) (*models.Metrics, error)
	GetAll(ctx context.Context) ([]*models.Metrics, error)
	Ping(ctx context.Context) error
}

type MemStorage struct {
	metrics map[string]*models.Metrics
	db      *sql.DB
}

func NewMemStorage(db *sql.DB) Storage {
	return &MemStorage{
		metrics: make(map[string]*models.Metrics),
		db:      db,
	}
}

func (m *MemStorage) Create(ctx context.Context, metric *models.Metrics) error {
	if _, err := m.Get(ctx, metric.ID); err == nil {
		return fmt.Errorf("metric %v already exists", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Update(ctx context.Context, metric *models.Metrics) error {
	if _, err := m.Get(ctx, metric.ID); err != nil {
		return fmt.Errorf("metric %v does not exist", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Get(ctx context.Context, ID string) (*models.Metrics, error) {
	if metric, ok := m.metrics[ID]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v does not exist", ID)
}

func (m *MemStorage) GetAll(ctx context.Context) ([]*models.Metrics, error) {
	var metrics []*models.Metrics
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *MemStorage) Ping(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("db is not configured")
	}
	return nil
}
