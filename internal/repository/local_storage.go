package repository

import (
	"fmt"
	models "github.com/krtech-it/metricagent/internal/model"
)

type MemStorage struct {
	metrics map[string]*models.Metrics
}

func NewMemStorage() models.Storage {
	return &MemStorage{
		metrics: make(map[string]*models.Metrics),
	}
}

func (m *MemStorage) Create(metric *models.Metrics) error {
	if _, err := m.Get(metric.ID); err == nil {
		return fmt.Errorf("metric %v already exists", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Update(metric *models.Metrics) error {
	if _, err := m.Get(metric.ID); err != nil {
		return fmt.Errorf("metric %v does not exist", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Get(ID string) (*models.Metrics, error) {
	if metric, ok := m.metrics[ID]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v does not exist", ID)
}
