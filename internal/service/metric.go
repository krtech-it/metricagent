package service

import (
	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
)

type MetricUseCase struct {
	storage repository.Storage
}

func NewMetricUseCase(storage repository.Storage) *MetricUseCase {
	return &MetricUseCase{
		storage: storage,
	}
}

func (m *MetricUseCase) Update(metric *models.Metrics) error {
	if metric.MType == models.Counter {
		oldMetric, err := m.storage.Get(metric.ID)
		if err != nil {
			return m.storage.Create(metric)
		} else {
			*metric.Delta += *oldMetric.Delta
			return m.storage.Update(metric)
		}
	} else if metric.MType == models.Gauge {
		if m.storage.Update(metric) != nil {
			return m.storage.Create(metric)
		}
	}
	return nil
}

func (m *MetricUseCase) GetMetric(ID string) (*models.Metrics, error) {
	return m.storage.Get(ID)
}

func (m *MetricUseCase) GetAllMetrics() ([]*models.Metrics, error) {
	return m.storage.GetAll()
}
