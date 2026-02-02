package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krtech-it/metricagent/internal/handler"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"log"
	"strconv"
)

func main() {
	addr := new(SetServer)
	if err := addr.Set(); err != nil {
		log.Fatal(err)
		return
	}
	storage := repository.NewMemStorage()
	metricUseCase := service.NewMetricUseCase(storage)

	h := handler.NewHandler(metricUseCase)

	r := gin.Default()
	r.LoadHTMLGlob("internal/templates/*")
	r.POST("/update/:metricType/:ID/:value", gin.WrapF(h.UpdateMetric))
	r.GET("/value/:metricType/:ID", h.GetMetric)
	r.GET("/", h.GetMainHTML)

	log.Println("Listening on port ", strconv.Itoa(addr.port))
	err := r.Run(addr.host + ":" + strconv.Itoa(addr.port))
	if err != nil {
		log.Fatal(err)
	}
}
