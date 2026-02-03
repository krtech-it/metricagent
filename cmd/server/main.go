package main

import (
	delivery "github.com/krtech-it/metricagent/internal/delivery/http"
	"log"
	"strconv"
)

func main() {
	addr, err := NewSetServer()
	if err != nil {
		log.Fatal(err)
	}
	router := delivery.NewRouter()
	log.Println("Listening on port ", strconv.Itoa(addr.port))
	err = router.Run(addr.host + ":" + strconv.Itoa(addr.port))
	if err != nil {
		log.Fatal(err)
	}
}
