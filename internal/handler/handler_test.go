package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandler() (*Handler, models.Storage) {
	storage := repository.NewMemStorage()
	metricUseCase := service.NewMetricUseCase(storage)
	return NewHandler(metricUseCase), storage
}

func TestUpdateMetricGaugeOK(t *testing.T) {
	h, storage := newTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/123.5", nil)
	rec := httptest.NewRecorder()

	h.UpdateMetric(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))

	metric, err := storage.Get("Alloc")
	require.NoError(t, err)
	require.NotNil(t, metric.Value)
	assert.InEpsilon(t, 123.5, *metric.Value, 0.0001)
}

func TestUpdateMetricCounterAccumulates(t *testing.T) {
	h, storage := newTestHandler()

	req1 := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/5", nil)
	rec1 := httptest.NewRecorder()
	h.UpdateMetric(rec1, req1)
	res1 := rec1.Result()
	defer res1.Body.Close()

	assert.Equal(t, http.StatusOK, res1.StatusCode)

	req2 := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/7", nil)
	rec2 := httptest.NewRecorder()
	h.UpdateMetric(rec2, req2)
	res2 := rec2.Result()
	defer res2.Body.Close()

	assert.Equal(t, http.StatusOK, res2.StatusCode)

	metric, err := storage.Get("PollCount")
	require.NoError(t, err)
	require.NotNil(t, metric.Delta)
	assert.Equal(t, int64(12), *metric.Delta)
}

func TestUpdateMetricErrors(t *testing.T) {
	h, _ := newTestHandler()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "wrong method",
			method:     http.MethodGet,
			path:       "/update/gauge/Alloc/1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid type",
			method:     http.MethodPost,
			path:       "/update/unknown/Alloc/1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid gauge value",
			method:     http.MethodPost,
			path:       "/update/gauge/Alloc/not-a-number",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid counter value",
			method:     http.MethodPost,
			path:       "/update/counter/PollCount/not-a-number",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing name",
			method:     http.MethodPost,
			path:       "/update/gauge//1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid path",
			method:     http.MethodPost,
			path:       "/update/gauge/Alloc",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			h.UpdateMetric(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
