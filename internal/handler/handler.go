package handler

import (
	"github.com/gin-gonic/gin"
	models "github.com/krtech-it/metricagent/internal/model"
	"github.com/krtech-it/metricagent/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	metricUseCase *service.MetricUseCase
}

func NewHandler(metricUseCase *service.MetricUseCase) *Handler {
	return &Handler{
		metricUseCase: metricUseCase,
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

	var metric *models.Metrics

	if metricType == models.Counter {
		delta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric = &models.Metrics{
			ID:    ID,
			MType: metricType,
			Delta: &delta,
			Value: nil,
			Hash:  "",
		}
	} else if metricType == models.Gauge {
		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		metric = &models.Metrics{
			ID:    ID,
			MType: metricType,
			Value: &valueFloat64,
			Delta: nil,
			Hash:  "",
		}
	}
	if err := h.metricUseCase.Update(metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetric(c *gin.Context) {
	metricType := c.Param("metricType")
	ID := c.Param("ID")
	if metricType != models.Gauge && metricType != models.Counter {
		c.String(http.StatusNotFound, "invalid path")
	}
	metric, err := h.metricUseCase.GetMetric(ID)
	if err != nil {
		c.String(http.StatusNotFound, "ID: %s does not exist", ID)
		return
	}
	if metric.MType != metricType {
		c.String(http.StatusBadRequest, "invalid path")
		return
	}
	switch metric.MType {
	case models.Gauge:
		s := strconv.FormatFloat(*metric.Value, 'f', 3, 64) // всегда 3 знака: "45.400"
		if strings.Contains(s, ".") {
			s = strings.TrimRight(strings.TrimRight(s, "0"), ".") // убираем нули и точку
		}
		c.String(http.StatusOK, s)
	case models.Counter:
		c.String(http.StatusOK, "%d", *metric.Delta)
	}
}

func (h *Handler) GetMainHTML(c *gin.Context) {
	metrics, err := h.metricUseCase.GetAllMetrics()
	if err != nil {
		log.Printf("handler: GetMainHTML, error: %s \n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "main_server.html", metrics)
}
