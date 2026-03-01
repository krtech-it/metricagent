package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host            string
	Port            int
	ReportInterval  int
	PollInterval    int
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	address := getEnv("ADDRESS", FlagServer.addr)
	if storeInterval, err := strconv.Atoi(getEnv("STORE_INTERVAL", strconv.Itoa(FlagServer.storeInterval))); err == nil {
		cfg.StoreInterval = storeInterval
	} else {
		cfg.StoreInterval = FlagServer.storeInterval
	}
	cfg.FileStoragePath = getEnv("FILE_STORAGE_PATH", FlagServer.fileStoragePath)
	if restore, err := strconv.ParseBool(getEnv("RESTORE", strconv.FormatBool(FlagServer.restore))); err == nil {
		cfg.Restore = restore
	} else {
		cfg.Restore = FlagServer.restore
	}

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

	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
