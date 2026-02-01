package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/handler"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"log"
)

func main() {
	storage := repository.NewMemStorage()
	metricUseCase := service.NewMetricUseCase(storage)

	h := handler.NewHandler(metricUseCase)

	r := gin.Default()

	r.POST("/update/:metricType/:ID/:value", gin.WrapF(h.UpdateMetric))
	r.GET("/update/value/:metricType/:ID", h.GetMetric)

	log.Println("Listening on port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
