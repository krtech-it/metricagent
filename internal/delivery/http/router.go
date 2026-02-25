package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/handler"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	storage := repository.NewMemStorage()
	metricUseCase := service.NewMetricUseCase(storage)

	h := handler.NewHandler(metricUseCase)
	r.LoadHTMLGlob("internal/templates/*")
	r.POST("/update/:metricType/:ID/:value", gin.WrapF(h.UpdateMetric))
	r.GET("/value/:metricType/:ID", h.GetMetric)
	r.GET("/", h.GetMainHTML)

	return r
}
