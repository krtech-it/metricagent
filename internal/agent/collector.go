package agent

import (
	"context"
	"math/rand"
	"runtime"
	"time"
)

type Collector struct {
	Storage map[string]interface{}
	counter int64
}

func NewCollector() *Collector {
	return &Collector{
		Storage: make(map[string]interface{}),
		counter: 0,
	}
}

func (c *Collector) Add(ctx context.Context, timeTicker time.Duration) {
	ticker := time.NewTicker(timeTicker * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			c.Storage["Alloc"] = memStats.Alloc
			c.Storage["BuckHashSys"] = memStats.BuckHashSys
			c.Storage["Frees"] = memStats.Frees
			c.Storage["GCCPUFraction"] = memStats.GCCPUFraction
			c.Storage["GCSys"] = memStats.GCSys
			c.Storage["HeapAlloc"] = memStats.HeapAlloc
			c.Storage["HeapIdle"] = memStats.HeapIdle
			c.Storage["HeapInuse"] = memStats.HeapInuse
			c.Storage["HeapObjects"] = memStats.HeapObjects
			c.Storage["HeapReleased"] = memStats.HeapReleased
			c.Storage["HeapSys"] = memStats.HeapSys
			c.Storage["LastGC"] = memStats.LastGC
			c.Storage["Lookups"] = memStats.Lookups
			c.Storage["MCacheInuse"] = memStats.MCacheInuse
			c.Storage["MCacheSys"] = memStats.MCacheSys
			c.Storage["MSpanInuse"] = memStats.MSpanInuse
			c.Storage["MSpanSys"] = memStats.MSpanSys
			c.Storage["Mallocs"] = memStats.Mallocs
			c.Storage["NextGC"] = memStats.NextGC
			c.Storage["NumForcedGC"] = memStats.NumForcedGC
			c.Storage["NumGC"] = memStats.NumGC
			c.Storage["OtherSys"] = memStats.OtherSys
			c.Storage["PauseTotalNs"] = memStats.PauseTotalNs
			c.Storage["StackInuse"] = memStats.StackInuse
			c.Storage["StackSys"] = memStats.StackSys
			c.Storage["Sys"] = memStats.Sys
			c.Storage["TotalAlloc"] = memStats.TotalAlloc

			c.counter++
			c.Storage["PoolCount"] = c.counter
			c.Storage["RandomValue"] = rand.Float64()
		}
	}
}
