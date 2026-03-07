package repository

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorageCreateAndGet(t *testing.T) {
	storage := NewMemStorage(nil)

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
	storage := NewMemStorage(nil)

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
	storage := NewMemStorage(nil)

	value := 2.0
	metric := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	}

	assert.Error(t, storage.Update(metric))
}

func TestMemStorageUpdateOverwrites(t *testing.T) {
	storage := NewMemStorage(nil)

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

func TestMemStorageGetAll(t *testing.T) {
	storage := NewMemStorage(nil)

	gaugeValue := 1.0
	counterValue := int64(3)
	require.NoError(t, storage.Create(&models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &gaugeValue,
	}))
	require.NoError(t, storage.Create(&models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: &counterValue,
	}))

	all, err := storage.GetAll()
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestMemStoragePingNilDB(t *testing.T) {
	storage := NewMemStorage(nil)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	err := storage.Ping(ctx)
	assert.Error(t, err)
}
