package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandler() (*Handler, repository.Storage) {
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

func TestGetMetric(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		params     gin.Params
		metricObj  *models.Metrics
		wantStatus int
		Value      float64
		Delta      int64
		wantValue  string
	}{
		{name: "test1",
			path:       "/value/gauge/Alloc",
			params:     gin.Params{{Key: "metricType", Value: "gauge"}, {Key: "ID", Value: "Alloc"}},
			wantValue:  "12.2",
			wantStatus: http.StatusOK,
			Value:      12.20,
			Delta:      0,
			metricObj: &models.Metrics{
				ID:    "Alloc",
				MType: models.Gauge,
			},
		},
		{name: "test2",
			path:       "/value/counter/PollCount",
			params:     gin.Params{{Key: "metricType", Value: "counter"}, {Key: "ID", Value: "PollCount"}},
			wantValue:  "10",
			wantStatus: http.StatusOK,
			Value:      0,
			Delta:      10,
			metricObj: &models.Metrics{
				ID:    "PollCount",
				MType: models.Counter,
			},
		},
		{name: "test3",
			path:       "/value/counter/PollCount",
			params:     gin.Params{{Key: "metricType", Value: "counter"}, {Key: "ID", Value: "PollCount"}},
			wantValue:  "ID: PollCount does not exist",
			wantStatus: http.StatusNotFound,
			Value:      0,
			Delta:      10,
			metricObj: &models.Metrics{
				ID:    "NotFound",
				MType: models.Counter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h, storage := newTestHandler()
			tt.metricObj.Value = &tt.Value
			tt.metricObj.Delta = &tt.Delta
			err := storage.Create(tt.metricObj)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Params = tt.params

			c.Request = httptest.NewRequest(http.MethodGet, tt.path, nil)

			h.GetMetric(c)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantValue, rec.Body.String())
		})
	}
}

func TestGetMainHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, storage := newTestHandler()

	value := 1.5
	err := storage.Create(&models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: &value,
	})
	require.NoError(t, err)

	router := gin.New()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", h.GetMainHTML)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Alloc")
}
