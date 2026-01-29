package main

import (
	"context"
	"fmt"
	"github.com/krtech-it/metricagent/internal/agent"
	"time"
)

func main() {
	collector := agent.NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	go collector.Add(ctx, 2000)
	defer cancel()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for name, value := range collector.Storage {
				err := agent.SendMetric(name, value)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	select {}
}
