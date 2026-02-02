package main

import (
	"context"
	"fmt"
	"github.com/krtech-it/metricagent/internal/agent"
	"log"
	"strconv"
	"time"
)

func main() {
	setAgent := new(SetAgent)
	if err := setAgent.Set(); err != nil {
		log.Fatal(err)
		return
	}
	collector := agent.NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	go collector.Add(ctx, time.Duration(setAgent.pollInterval*1000))
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Duration(setAgent.reportInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for name, value := range collector.CopyStorage() {
				err := agent.SendMetric(name, value, setAgent.host+":"+strconv.Itoa(setAgent.port))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	select {}
}
