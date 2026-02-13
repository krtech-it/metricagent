package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host           string
	Port           int
	ReportInterval int
	PoolInterval   int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	address := getEnv("ADDRESS", FlagAgent.addr)
	pollStr := getEnv("POLL_INTERVAL", FlagAgent.pollInterval)
	reportStr := getEnv("REPORT_INTERVAL", FlagAgent.reportInterval)

	args := strings.Split(address, ":")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid server address %s", address)
	}
	cfg.Host = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("port is not int: %w", err)
	}
	cfg.Port = port
	pollInt, err := strconv.Atoi(pollStr)
	if err != nil {
		return nil, fmt.Errorf("poll interval is not int: %w", err)
	}
	cfg.PoolInterval = pollInt
	reportInt, err := strconv.Atoi(reportStr)
	if err != nil {
		return nil, fmt.Errorf("report interval is not int: %w", err)
	}
	cfg.ReportInterval = reportInt
	if cfg.PoolInterval == 0 || cfg.ReportInterval == 0 {
		return nil, fmt.Errorf("report interval and report interval must be greater than zero")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
