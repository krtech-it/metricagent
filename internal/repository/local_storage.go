package repository

import (
	"context"
	"database/sql"
	"fmt"
	models "github.com/krtech-it/metricagent/internal/model"
	"sync"
)

type Storage interface {
	Update(ctx context.Context, metric *models.Metrics) error
	Create(ctx context.Context, metric *models.Metrics) error
	Upsert(ctx context.Context, metrics []*models.Metrics) error
	Get(ctx context.Context, id string) (*models.Metrics, error)
	GetAll(ctx context.Context) ([]*models.Metrics, error)
	Ping(ctx context.Context) error
}

type MemStorage struct {
	mu      sync.Mutex
	metrics map[string]*models.Metrics
	db      *sql.DB
}

func NewMemStorage(db *sql.DB) Storage {
	return &MemStorage{
		metrics: make(map[string]*models.Metrics),
		db:      db,
		mu:      sync.Mutex{},
	}
}

func (m *MemStorage) getLocked(ID string) (*models.Metrics, error) {
	if metric, ok := m.metrics[ID]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v does not exist", ID)
}

func (m *MemStorage) createLocked(metric *models.Metrics) error {
	if _, err := m.getLocked(metric.ID); err == nil {
		return fmt.Errorf("metric %v already exists", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) updateLocked(metric *models.Metrics) error {
	if _, err := m.getLocked(metric.ID); err != nil {
		return fmt.Errorf("metric %v does not exist", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Create(ctx context.Context, metric *models.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.createLocked(metric)
}

func (m *MemStorage) Update(ctx context.Context, metric *models.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.updateLocked(metric)
}

func (m *MemStorage) Upsert(ctx context.Context, metrics []*models.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, metric := range metrics {
		oldMetric, err := m.getLocked(metric.ID)
		if err != nil {
			m.createLocked(metric)
			continue
		}
		if metric.MType == models.Counter {
			if metric.Delta != nil && oldMetric.Delta != nil {
				metricDelta := *oldMetric.Delta + *metric.Delta
				metric.Delta = &metricDelta
			} else if oldMetric.Delta != nil {
				metric.Delta = oldMetric.Delta
			}
		}
		m.updateLocked(metric)
	}
	return nil
}

func (m *MemStorage) Get(ctx context.Context, ID string) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.getLocked(ID)

}

func (m *MemStorage) GetAll(ctx context.Context) ([]*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
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
