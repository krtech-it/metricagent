package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	models "github.com/krtech-it/metricagent/internal/agent/dto"
	"net"
	"net/http"
	"slices"
	"strconv"
	"syscall"
	"time"
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

func SendMetricJSON(name string, value interface{}, host string) error {
	var mType string
	var requestMetric models.RequestMetricUpdate
	if slices.Contains(gaugeArea[:], name) {
		mType = "gauge"
		switch v := value.(type) {
		case float64:
			requestMetric.Value = &v
		case uint64:
			v64 := float64(v)
			requestMetric.Value = &v64
		case uint32:
			v64 := float64(v)
			requestMetric.Value = &v64
		case int64:
			v64 := float64(v)
			requestMetric.Value = &v64
		}
	} else {
		mType = "counter"
		v, _ := value.(int64)
		requestMetric.Delta = &v
	}
	requestMetric.MType = mType
	requestMetric.ID = name
	url := fmt.Sprintf("http://%s/update/", host)
	body, err := json.Marshal(requestMetric)
	if err != nil {
		return err
	}
	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	if _, err := gz.Write(body); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, &gzBuf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
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

func SendMetricsJSON(items map[string]interface{}, host string) error {
	const (
		maxRetries = 3
		baseDelay  = 2 * time.Second
	)

	var lastErr error
	delay := 1 * time.Second
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := SendMetricsOnce(items, host)
		if err == nil {
			return nil
		}
		lastErr = err
		if !IsRetriableError(err) {
			return fmt.Errorf("non-retriable error - %d: %w", attempt, err)
		}
		if attempt != 1 && attempt != maxRetries {
			delay += baseDelay
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func SendMetricsOnce(items map[string]interface{}, host string) error {
	var mType string
	var requestMetric models.RequestMetricUpdate
	var requestMetrics []models.RequestMetricUpdate
	for name, value := range items {
		if slices.Contains(gaugeArea[:], name) {
			mType = "gauge"
			switch v := value.(type) {
			case float64:
				requestMetric.Value = &v
			case uint64:
				v64 := float64(v)
				requestMetric.Value = &v64
			case uint32:
				v64 := float64(v)
				requestMetric.Value = &v64
			case int64:
				v64 := float64(v)
				requestMetric.Value = &v64
			}
		} else {
			mType = "counter"
			v, ok := value.(int64)
			if !ok {
				return fmt.Errorf("unexpected counter value type for %s: %T", name, value)
			}
			requestMetric.Delta = &v
		}
		requestMetric.MType = mType
		requestMetric.ID = name
		requestMetrics = append(requestMetrics, requestMetric)
	}

	url := fmt.Sprintf("http://%s/updates/", host)
	body, err := json.Marshal(requestMetrics)
	if err != nil {
		return fmt.Errorf("failed to marshal request metrics: %w", err)
	}
	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	if _, err := gz.Write(body); err != nil {
		return fmt.Errorf("failed to gzip request body: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("gzip close faile: %w", err)
	}
	req, err := http.NewRequest("POST", url, &gzBuf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %s: %w", resp.Status, errors.New(resp.Status))
		}
		return fmt.Errorf("bad status: %s, url - %s", resp.Status, url)
	}
	return nil
}
