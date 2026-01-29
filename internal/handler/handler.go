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
	path_args := strings.Split(path, "/")
	if len(path_args) != 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	metric_type := path_args[0]
	ID := path_args[1]
	value := path_args[2]
	if !(metric_type == models.Gauge || metric_type == models.Counter) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	if ID == "" {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	}
	if metric_type == models.Counter {
		delta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric := &models.Metrics{
			ID:    ID,
			MType: metric_type,
			Delta: &delta,
			Value: nil,
			Hash:  "",
		}
		old_metric, err := h.storage.Get(metric.ID)
		if err != nil {
			h.storage.Create(metric)
		} else {
			delta += *old_metric.Delta
			h.storage.Update(metric)
		}
	} else if metric_type == models.Gauge {
		value_float64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric := &models.Metrics{
			ID:    ID,
			MType: metric_type,
			Value: &value_float64,
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
