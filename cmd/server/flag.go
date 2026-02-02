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

func (s *SetServer) Set() error {
	var addr string
	flag.StringVar(&addr, "a", "localhost:8080", "server listen address")
	flag.Parse()
	args := strings.Split(addr, ":")
	if len(args) != 2 {
		return errors.New("invalid server address")
	}
	s.host = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	s.port = port
	return nil
}
