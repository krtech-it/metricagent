package agent

import (
	"fmt"
	"net/http"
	"slices"
)

var (
	gaugeArea = [...]string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
	}
)

func SendMetric(name string, value interface{}) error {
	var mType string
	if slices.Contains(gaugeArea[:], name) {
		mType = "gauge"
	} else {
		mType = "counter"
	}
	var sendValue string
	switch value.(type) {
	case uint64:
		sendValue = fmt.Sprintf("%d", value)
	case float64:
		sendValue = fmt.Sprintf("%f", value)
	case uint32:
		sendValue = fmt.Sprintf("%d", value)
	case int64:
		sendValue = fmt.Sprintf("%d", value)
	}
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", mType, name, sendValue)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s, url - %s", resp.Status, url)
	}
	return nil
}
