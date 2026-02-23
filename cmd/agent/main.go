package main

import (
	"github.com/krtech-it/metricagent/internal/agent"
	"github.com/krtech-it/metricagent/internal/agent/config"
	"log"
	"strconv"
	"time"
)

func main() {
	config.ParseFlags()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	collector := agent.NewCollector()
	tickerPool := time.NewTicker(time.Duration(cfg.PoolInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	done := make(chan bool)

	defer tickerPool.Stop()
	defer tickerReport.Stop()
	defer close(done)

	for {
		select {
		case <-done:
			return
		case <-tickerPool.C:
			collector.Add()
		case <-tickerReport.C:
			errFlag := false
			for name, value := range collector.CopyStorage() {
				err := agent.SendMetricJSON(name, value, cfg.Host+":"+strconv.Itoa(cfg.Port))
				if err != nil {
					log.Printf("error send metric: %s \n", err)
					if name == "PoollCount" {
						errFlag = true
						break
					}
				}
			}
			if !errFlag {
				collector.ResetCounter()
			}
		}
	}
}
