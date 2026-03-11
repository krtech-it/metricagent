package backuper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/krtech-it/metricagent/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBackuperWriteRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	backup, err := NewBackuper(path, zap.NewNop())
	require.NoError(t, err)

	gaugeValue := 12.5
	counterValue := int64(7)
	input := []*dto.ResponseGetMetric{
		{
			MainMetric: dto.MainMetric{ID: "Alloc", MType: "gauge"},
			Value:      &gaugeValue,
		},
		{
			MainMetric: dto.MainMetric{ID: "PollCount", MType: "counter"},
			Delta:      &counterValue,
		},
	}

	require.NoError(t, backup.WriteEvent(input))

	result, err := backup.ReadEvent()
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Alloc", result[0].ID)
	require.NotNil(t, result[0].Value)
	assert.InEpsilon(t, 12.5, *result[0].Value, 0.0001)
	assert.Equal(t, "PollCount", result[1].ID)
	require.NotNil(t, result[1].Delta)
	assert.Equal(t, int64(7), *result[1].Delta)
}

func TestBackuperOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	backup, err := NewBackuper(path, zap.NewNop())
	require.NoError(t, err)

	first := int64(1)
	second := int64(2)
	require.NoError(t, backup.WriteEvent([]*dto.ResponseGetMetric{
		{MainMetric: dto.MainMetric{ID: "PollCount", MType: "counter"}, Delta: &first},
	}))
	require.NoError(t, backup.WriteEvent([]*dto.ResponseGetMetric{
		{MainMetric: dto.MainMetric{ID: "PollCount", MType: "counter"}, Delta: &second},
	}))

	result, err := backup.ReadEvent()
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].Delta)
	assert.Equal(t, int64(2), *result[0].Delta)
}

func TestBackuperReadEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	backup, err := NewBackuper(path, zap.NewNop())
	require.NoError(t, err)

	_, err = backup.ReadEvent()
	assert.Error(t, err)
}

func TestBackuperReadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	require.NoError(t, os.WriteFile(path, []byte("{not-json"), 0o644))

	backup, err := NewBackuper(path, zap.NewNop())
	require.NoError(t, err)

	_, err = backup.ReadEvent()
	assert.Error(t, err)
}

func TestBackuperInvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "metrics.json")
	_, err := NewBackuper(path, zap.NewNop())
	assert.Error(t, err)
}
