package agent

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
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

func SendMetric(name string, value interface{}, host string) error {
	var mType string
	if slices.Contains(gaugeArea[:], name) {
		mType = "gauge"
	} else {
		mType = "counter"
	}
	var sendValue string
	switch v := value.(type) {
	case uint64:
		sendValue = strconv.Itoa(int(v))
	case float64:
		sendValue = strconv.FormatFloat(v, 'f', -1, 64)
	case uint32:
		sendValue = strconv.Itoa(int(v))
	case int64:
		sendValue = strconv.Itoa(int(v))
	default:
		return fmt.Errorf("ошибка создания url невалидное значение")
	}
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, mType, name, sendValue)
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
