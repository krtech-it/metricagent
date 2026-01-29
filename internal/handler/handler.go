package handler

import (
	models "github.com/krtech-it/metricagent/internal/model"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	storage models.Storage
}

func NewHandler(storage models.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/")
	pathArgs := strings.Split(path, "/")
	if len(pathArgs) != 3 {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	}
	metricType := pathArgs[0]
	ID := pathArgs[1]
	value := pathArgs[2]
	if !(metricType == models.Gauge || metricType == models.Counter) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	if ID == "" {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	}
	if metricType == models.Counter {
		delta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric := &models.Metrics{
			ID:    ID,
			MType: metricType,
			Delta: &delta,
			Value: nil,
			Hash:  "",
		}
		oldMetric, err := h.storage.Get(metric.ID)
		if err != nil {
			h.storage.Create(metric)
		} else {
			delta += *oldMetric.Delta
			h.storage.Update(metric)
		}
	} else if metricType == models.Gauge {
		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric := &models.Metrics{
			ID:    ID,
			MType: metricType,
			Value: &valueFloat64,
			Delta: nil,
			Hash:  "",
		}
		if h.storage.Update(metric) != nil {
			h.storage.Create(metric)
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
