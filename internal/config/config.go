package config

import (
	"flag"
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
	var (
		addr = flag.String("a", "localhost:8080", "server listen address")
	)
	flag.Parse()
	address := getEnv("ADDRESS", *addr)

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
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
