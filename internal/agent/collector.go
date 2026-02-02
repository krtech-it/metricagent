package agent

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Collector struct {
	mu      sync.RWMutex
	storage map[string]interface{}
	counter int64
}

func NewCollector() *Collector {
	return &Collector{
		storage: make(map[string]interface{}),
		counter: 0,
		mu:      sync.RWMutex{},
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
			c.mu.Lock()
			c.storage["Alloc"] = memStats.Alloc
			c.storage["BuckHashSys"] = memStats.BuckHashSys
			c.storage["Frees"] = memStats.Frees
			c.storage["GCCPUFraction"] = memStats.GCCPUFraction
			c.storage["GCSys"] = memStats.GCSys
			c.storage["HeapAlloc"] = memStats.HeapAlloc
			c.storage["HeapIdle"] = memStats.HeapIdle
			c.storage["HeapInuse"] = memStats.HeapInuse
			c.storage["HeapObjects"] = memStats.HeapObjects
			c.storage["HeapReleased"] = memStats.HeapReleased
			c.storage["HeapSys"] = memStats.HeapSys
			c.storage["LastGC"] = memStats.LastGC
			c.storage["Lookups"] = memStats.Lookups
			c.storage["MCacheInuse"] = memStats.MCacheInuse
			c.storage["MCacheSys"] = memStats.MCacheSys
			c.storage["MSpanInuse"] = memStats.MSpanInuse
			c.storage["MSpanSys"] = memStats.MSpanSys
			c.storage["Mallocs"] = memStats.Mallocs
			c.storage["NextGC"] = memStats.NextGC
			c.storage["NumForcedGC"] = memStats.NumForcedGC
			c.storage["NumGC"] = memStats.NumGC
			c.storage["OtherSys"] = memStats.OtherSys
			c.storage["PauseTotalNs"] = memStats.PauseTotalNs
			c.storage["StackInuse"] = memStats.StackInuse
			c.storage["StackSys"] = memStats.StackSys
			c.storage["Sys"] = memStats.Sys
			c.storage["TotalAlloc"] = memStats.TotalAlloc

			c.counter++
			c.storage["PoolCount"] = c.counter
			c.storage["RandomValue"] = rand.Float64()
			c.mu.Unlock()
		}
	}
}

func (c *Collector) CopyStorage() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	storageCopy := make(map[string]interface{})
	for k, v := range c.storage {
		storageCopy[k] = v
	}
	return storageCopy
}
