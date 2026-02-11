package service

import (
	"testing"

	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUseCase_UpdateCounterAccumulates(t *testing.T) {
	storage := repository.NewMemStorage()
	useCase := NewMetricUseCase(storage)

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
	storage := repository.NewMemStorage()
	useCase := NewMetricUseCase(storage)

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
	storage := repository.NewMemStorage()
	useCase := NewMetricUseCase(storage)

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
