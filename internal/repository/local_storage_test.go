package repository

import (
	"testing"

	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorageCreateAndGet(t *testing.T) {
	storage := NewMemStorage()

	value := 1.5
	metric := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	}

	err := storage.Create(metric)
	require.NoError(t, err)

	got, err := storage.Get("Alloc")
	require.NoError(t, err)
	require.NotNil(t, got.Value)
	assert.InEpsilon(t, 1.5, *got.Value, 0.0001)
}

func TestMemStorageCreateExisting(t *testing.T) {
	storage := NewMemStorage()

	value := 1.5
	metric := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	}

	require.NoError(t, storage.Create(metric))
	assert.Error(t, storage.Create(metric))
}

func TestMemStorageUpdateMissing(t *testing.T) {
	storage := NewMemStorage()

	value := 2.0
	metric := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	}

	assert.Error(t, storage.Update(metric))
}

func TestMemStorageUpdateOverwrites(t *testing.T) {
	storage := NewMemStorage()

	first := 1.0
	metric := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &first,
	}
	require.NoError(t, storage.Create(metric))

	second := 2.0
	updated := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &second,
	}
	require.NoError(t, storage.Update(updated))

	got, err := storage.Get("Alloc")
	require.NoError(t, err)
	require.NotNil(t, got.Value)
	assert.InEpsilon(t, 2.0, *got.Value, 0.0001)
}
