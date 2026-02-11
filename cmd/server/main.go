package main

import (
	"github.com/krtech-it/metricagent/internal/config"
	delivery "github.com/krtech-it/metricagent/internal/delivery/http"
	"log"
	"strconv"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	router := delivery.NewRouter()
	log.Println("Listening on port ", strconv.Itoa(cfg.Port))
	err = router.Run(cfg.Host + ":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
