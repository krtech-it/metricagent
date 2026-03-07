package service

import (
	"context"
	"github.com/krtech-it/metricagent/internal/backuper"
	"github.com/krtech-it/metricagent/internal/config"
	"github.com/krtech-it/metricagent/internal/delivery/http/dto"
	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
)

type MetricUseCase struct {
	storage    repository.Storage
	backup     backuper.BackupInterface
	cfg        *config.Config
	flagBackup bool
}

func NewMetricUseCase(storage repository.Storage, backup backuper.BackupInterface, cfg *config.Config) *MetricUseCase {
	return &MetricUseCase{
		storage:    storage,
		backup:     backup,
		cfg:        cfg,
		flagBackup: true,
	}
}

func (m *MetricUseCase) Update(ctx context.Context, metric *models.Metrics) error {
	var errResult error
	if metric.MType == models.Counter {
		oldMetric, err := m.storage.Get(ctx, metric.ID)
		if err != nil {
			errResult = m.storage.Create(ctx, metric)
		} else {
			*metric.Delta += *oldMetric.Delta
			errResult = m.storage.Update(ctx, metric)
		}
	} else if metric.MType == models.Gauge {
		if m.storage.Update(ctx, metric) != nil {
			errResult = m.storage.Create(ctx, metric)
		}
	}
	if m.flagBackup && m.cfg.StoreInterval == 0 && m.cfg.TypeDB == "file" {
		if err := m.WriteBackupAllMetrics(ctx); err != nil {
			return err
		}
	}
	return errResult
}

func (m *MetricUseCase) GetMetric(ctx context.Context, ID string) (*models.Metrics, error) {
	return m.storage.Get(ctx, ID)
}

func (m *MetricUseCase) GetAllMetrics(ctx context.Context) ([]*models.Metrics, error) {
	return m.storage.GetAll(ctx)
}

func (m *MetricUseCase) WriteBackupAllMetrics(ctx context.Context) error {
	allMetrics, err := m.storage.GetAll(ctx)
	if err != nil {
		return err
	}
	var metrics []*dto.ResponseGetMetric
	for _, metric := range allMetrics {
		metricDto := &dto.ResponseGetMetric{
			MainMetric: dto.MainMetric{
				MType: metric.MType,
				ID:    metric.ID,
			},
		}
		if metric.MType == models.Counter {
			metricDto.Delta = metric.Delta
		} else if metric.MType == models.Gauge {
			metricDto.Value = metric.Value
		}
		metrics = append(metrics, metricDto)
	}
	return m.backup.WriteEvent(metrics)
}

func (m *MetricUseCase) ReadBackupAllMetrics(ctx context.Context) error {
	m.flagBackup = false
	allMetrics, err := m.backup.ReadEvent()
	if err != nil {
		return err
	}
	for _, metric := range allMetrics {
		metricStorage := models.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
			Hash:  "",
		}
		if metric.MType == models.Counter {
			metricStorage.Delta = metric.Delta
			metricStorage.Value = nil
		} else if metric.MType == models.Gauge {
			metricStorage.Value = metric.Value
			metricStorage.Delta = nil
		}
		_ = m.Update(ctx, &metricStorage)
	}
	m.flagBackup = true
	return nil
}

func (m *MetricUseCase) Ping(ctx context.Context) error {
	return m.storage.Ping(ctx)
}
