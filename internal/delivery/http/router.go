package delivery

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/backuper"
	"github.com/krtech-it/metricagent/internal/config"
	"github.com/krtech-it/metricagent/internal/handler"
	"github.com/krtech-it/metricagent/internal/middleware"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"go.uber.org/zap"
	"time"
)

func NewRouter(logger *zap.Logger, cfg *config.Config, db *sql.DB) *gin.Engine {
	r := gin.Default()
	var storage repository.Storage
	if db != nil {
		storage = repository.NewDBStorage(db)
	} else {
		storage = repository.NewMemStorage(nil)
	}

	switch {
	case db != nil:
		cfg.TypeDB = "postgres"
	case cfg.FileStoragePath != "":
		cfg.TypeDB = "file"
	default:
		cfg.TypeDB = "memory"
	}

	backupService, err := backuper.NewBackuper(cfg.FileStoragePath, logger)
	if err != nil {
		logger.Error("failed to init backup service", zap.Error(err))
	}
	metricUseCase := service.NewMetricUseCase(storage, backupService, cfg)

	if cfg.TypeDB != "postgres" {
		if cfg.TypeDB == "file" {

			if cfg.Restore {
				if err := metricUseCase.ReadBackupAllMetrics(context.Background()); err != nil {
					logger.Warn("ReadBackupAllMetrics failed", zap.Error(err))
				}
			}
			if db == nil && cfg.StoreInterval > 0 {
				ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
				go func() {
					for range ticker.C {
						if err := metricUseCase.WriteBackupAllMetrics(context.Background()); err != nil {
							logger.Error(err.Error())
						}
					}
				}()
			}
		}
	}

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.DecompressMiddleware())
	r.Use(middleware.CheckHashMiddleware(cfg))
	r.Use(middleware.GzipMiddleware())
	r.Use(middleware.ResponseHashMiddleware(cfg))

	h := handler.NewHandler(metricUseCase, logger, cfg)
	r.LoadHTMLGlob("internal/templates/*")
	r.GET("/ping", h.Ping)
	r.POST("/update/", h.UpdateMetricJSON)
	r.POST("/updates/", h.UpdatesMetricJSON)
	r.POST("/update/:metricType/:ID/:value", gin.WrapF(h.UpdateMetric))
	r.POST("/value/", h.GetMetricJSON)
	r.GET("/value/:metricType/:ID", h.GetMetric)
	r.GET("/", h.GetMainHTML)

	return r
}
