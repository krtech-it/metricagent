package main

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type SetServer struct {
	host string
	port int
}

func NewSetServer() (*SetServer, error) {
	var addr string
	server := &SetServer{}
	flag.StringVar(&addr, "a", "localhost:8080", "server listen address")
	flag.Parse()
	args := strings.Split(addr, ":")
	if len(args) != 2 {
		return nil, errors.New("invalid server address")
	}
	server.host = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	server.port = port
	return server, nil
}
