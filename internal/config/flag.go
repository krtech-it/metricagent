package config

import (
	"flag"
)

type SetServer struct {
	addr string
}

var FlagServer = SetServer{}

func ParseFlags() {
	flag.StringVar(&FlagServer.addr, "a", "localhost:8080", "server listen address")
	if !flag.Parsed() {
		flag.Parse()
	}
}
