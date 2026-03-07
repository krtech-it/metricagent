package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	models "github.com/krtech-it/metricagent/internal/model"
	"time"
)

type Storage interface {
	Update(metric *models.Metrics) error
	Create(metric *models.Metrics) error
	Get(id string) (*models.Metrics, error)
	GetAll() ([]*models.Metrics, error)
	Ping(ctx *gin.Context) error
}

type MemStorage struct {
	metrics map[string]*models.Metrics
	db      *sql.DB
}

func NewMemStorage(db *sql.DB) Storage {
	return &MemStorage{
		metrics: make(map[string]*models.Metrics),
		db:      db,
	}
}

func (m *MemStorage) Create(metric *models.Metrics) error {
	if _, err := m.Get(metric.ID); err == nil {
		return fmt.Errorf("metric %v already exists", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Update(metric *models.Metrics) error {
	if _, err := m.Get(metric.ID); err != nil {
		return fmt.Errorf("metric %v does not exist", metric.ID)
	}
	m.metrics[metric.ID] = metric
	return nil
}

func (m *MemStorage) Get(ID string) (*models.Metrics, error) {
	if metric, ok := m.metrics[ID]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v does not exist", ID)
}

func (m *MemStorage) GetAll() ([]*models.Metrics, error) {
	var metrics []*models.Metrics
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *MemStorage) Ping(ctx *gin.Context) error {
	if m.db == nil {
		return fmt.Errorf("db is not configured")
	}
	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := m.db.PingContext(ctxPing); err != nil {
		return err
	}
	return nil
}
