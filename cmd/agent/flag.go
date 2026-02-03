package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type SetAgent struct {
	host           string
	port           int
	reportInterval int
	pollInterval   int
}

func NewSetAgent() (*SetAgent, error) {
	var addr string
	s := &SetAgent{}
	flag.StringVar(&addr, "a", "localhost:8080", "server listen address")
	flag.IntVar(&s.pollInterval, "p", 2, "poll interval seconds")
	flag.IntVar(&s.reportInterval, "r", 10, "report interval seconds")
	flag.Parse()
	args := strings.Split(addr, ":")
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid server address %s", addr))
	}
	s.host = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("port is not int: %w", err)
	}
	s.port = port
	return s, nil
}
