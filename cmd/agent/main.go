package main

import (
	"log"
	"strconv"
	"time"

	"github.com/krtech-it/metricagent/internal/agent"
	"github.com/krtech-it/metricagent/internal/agent/config"
)

func main() {
	config.ParseFlags()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	collector := agent.NewCollector()
	tickerPool := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	tickerPool2 := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	done := make(chan bool)
	jobs := make(chan map[string]interface{})

	defer tickerPool.Stop()
	defer tickerReport.Stop()
	defer close(jobs)
	defer close(done)

	goWorkers(jobs, cfg.RateLimit, cfg)

	go func() {
		select {
		case <-done:
			return
		case <-tickerPool2.C:
			collector.AddGopsutil()
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-tickerPool.C:
			collector.Add()
		case <-tickerReport.C:
			jobs <- collector.CopyStorage()
			//collector.ResetCounter()
		}
	}
}

func goWorkers(jobs chan map[string]interface{}, numWorkers int, cfg *config.Config) {
	for i := 0; i < numWorkers; i++ {
		go func() {
			for job := range jobs {
				err := agent.SendMetricsJSON(job, cfg.Host+":"+strconv.Itoa(cfg.Port), cfg)
				if err != nil {
					log.Printf("error send metric: %s \n", err)
				}
			}
		}()
	}
}
