package delivery

import (
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

func NewRouter(logger *zap.Logger, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	storage := repository.NewMemStorage()
	backupService, err := backuper.NewBackuper(cfg.FileStoragePath, logger)
	if err != nil {
		logger.Error("failed to init backup service", zap.Error(err))
	}
	metricUseCase := service.NewMetricUseCase(storage, backupService, cfg)

	if cfg.Restore {
		if err := metricUseCase.ReadBackupAllMetrics(); err != nil {
			logger.Warn("ReadBackupAllMetrics failed", zap.Error(err))
		}
	}
	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		go func() {
			for range ticker.C {
				if err := metricUseCase.WriteBackupAllMetrics(); err != nil {
					logger.Error(err.Error())
				}
			}
		}()
	}

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.DecompressMiddleware())
	r.Use(middleware.GzipMiddleware())

	h := handler.NewHandler(metricUseCase, logger, cfg)
	r.LoadHTMLGlob("internal/templates/*")
	r.POST("/update/", h.UpdateMetricJSON)
	r.POST("/update/:metricType/:ID/:value", gin.WrapF(h.UpdateMetric))
	r.POST("/value/", h.GetMetricJSON)
	r.GET("/value/:metricType/:ID", h.GetMetric)
	r.GET("/", h.GetMainHTML)

	return r
}
