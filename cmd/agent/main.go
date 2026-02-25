package main

import (
	"github.com/krtech-it/metricagent/internal/agent"
	"log"
	"strconv"
	"time"
)

func main() {
	setAgent, err := NewSetAgent()
	if err != nil {
		log.Fatal(err)
	}
	collector := agent.NewCollector()
	tickerPool := time.NewTicker(time.Duration(setAgent.pollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(setAgent.reportInterval) * time.Second)
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
				err := agent.SendMetric(name, value, setAgent.host+":"+strconv.Itoa(setAgent.port))
				if err != nil {
					log.Printf("error send metric: %s \n", err)
					if name == "PoolCount" {
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
