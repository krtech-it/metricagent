package service

import (
	"testing"

	"github.com/krtech-it/metricagent/internal/config"
	dto "github.com/krtech-it/metricagent/internal/delivery/http/dto"
	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeBackup struct {
	writeCalls int
	events     []*dto.ResponseGetMetric
}

func (f *fakeBackup) WriteEvent(metrics []*dto.ResponseGetMetric) error {
	f.writeCalls++
	f.events = metrics
	return nil
}

func (f *fakeBackup) ReadEvent() ([]*dto.ResponseGetMetric, error) {
	return f.events, nil
}

func TestMetricUseCase_UpdateCounterAccumulates(t *testing.T) {
	storage := repository.NewMemStorage(nil)
	cfg := &config.Config{StoreInterval: 1}
	useCase := NewMetricUseCase(storage, nil, cfg)

	first := int64(5)
	err := useCase.Update(&models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: &first,
	})
	require.NoError(t, err)

	second := int64(7)
	err = useCase.Update(&models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: &second,
	})
	require.NoError(t, err)

	metric, err := storage.Get("PollCount")
	require.NoError(t, err)
	require.NotNil(t, metric.Delta)
	assert.Equal(t, int64(12), *metric.Delta)
}

func TestMetricUseCase_UpdateGaugeReplaces(t *testing.T) {
	storage := repository.NewMemStorage(nil)
	cfg := &config.Config{StoreInterval: 1}
	useCase := NewMetricUseCase(storage, nil, cfg)

	first := 1.5
	err := useCase.Update(&models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &first,
	})
	require.NoError(t, err)

	second := 2.5
	err = useCase.Update(&models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &second,
	})
	require.NoError(t, err)

	metric, err := storage.Get("Alloc")
	require.NoError(t, err)
	require.NotNil(t, metric.Value)
	assert.InEpsilon(t, 2.5, *metric.Value, 0.0001)
}

func TestMetricUseCase_UpdateCreatesWhenMissing(t *testing.T) {
	storage := repository.NewMemStorage(nil)
	cfg := &config.Config{StoreInterval: 1}
	useCase := NewMetricUseCase(storage, nil, cfg)

	value := 10.0
	err := useCase.Update(&models.Metrics{
		ID:    "HeapAlloc",
		MType: models.Gauge,
		Value: &value,
	})
	require.NoError(t, err)

	metric, err := storage.Get("HeapAlloc")
	require.NoError(t, err)
	require.NotNil(t, metric.Value)
	assert.InEpsilon(t, 10.0, *metric.Value, 0.0001)
}

func TestMetricUseCase_UpdateWritesBackupWhenSync(t *testing.T) {
	storage := repository.NewMemStorage(nil)
	cfg := &config.Config{StoreInterval: 0}
	backup := &fakeBackup{}
	useCase := NewMetricUseCase(storage, backup, cfg)

	value := 1.0
	err := useCase.Update(&models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, backup.writeCalls)
}

func TestMetricUseCase_ReadBackupAllMetricsLoadsStorage(t *testing.T) {
	storage := repository.NewMemStorage(nil)
	cfg := &config.Config{StoreInterval: 0}
	backup := &fakeBackup{
		events: []*dto.ResponseGetMetric{
			{
				MainMetric: dto.MainMetric{ID: "Alloc", MType: models.Gauge},
				Value:      func() *float64 { v := 2.5; return &v }(),
			},
			{
				MainMetric: dto.MainMetric{ID: "PollCount", MType: models.Counter},
				Delta:      func() *int64 { v := int64(3); return &v }(),
			},
		},
	}
	useCase := NewMetricUseCase(storage, backup, cfg)

	require.NoError(t, useCase.ReadBackupAllMetrics())

	gauge, err := storage.Get("Alloc")
	require.NoError(t, err)
	require.NotNil(t, gauge.Value)
	assert.InEpsilon(t, 2.5, *gauge.Value, 0.0001)

	counter, err := storage.Get("PollCount")
	require.NoError(t, err)
	require.NotNil(t, counter.Delta)
	assert.Equal(t, int64(3), *counter.Delta)
}
