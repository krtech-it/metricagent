package main

import (
	delivery "github.com/krtech-it/metricagent/internal/delivery/http"
	"log"
	"strconv"
)

func main() {
	addr := new(SetServer)
	if err := addr.Set(); err != nil {
		log.Fatal(err)
		return
	}
	router := delivery.NewRouter()
	log.Println("Listening on port ", strconv.Itoa(addr.port))
	err := router.Run(addr.host + ":" + strconv.Itoa(addr.port))
	if err != nil {
		log.Fatal(err)
	}
}
