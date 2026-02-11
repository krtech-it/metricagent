package main

import (
	"github.com/krtech-it/metricagent/internal/handler"
	"github.com/krtech-it/metricagent/internal/repository"
	"github.com/krtech-it/metricagent/internal/service"
	"log"
	"net/http"
)

func main() {
	storage := repository.NewMemStorage()
	metricUseCase := service.NewMetricUseCase(storage)

	h := handler.NewHandler(metricUseCase)
	http.HandleFunc("/update/", h.UpdateMetric)

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
