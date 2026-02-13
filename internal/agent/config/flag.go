package config

import (
	"flag"
)

type SetAgent struct {
	addr           string
	reportInterval string
	pollInterval   string
}

var FlagAgent = SetAgent{}

func ParseFlags() {
	flag.StringVar(&FlagAgent.addr, "a", "localhost:8080", "server listen address")
	flag.StringVar(&FlagAgent.pollInterval, "p", "2", "poll interval seconds")
	flag.StringVar(&FlagAgent.reportInterval, "r", "10", "report interval seconds")
	if !flag.Parsed() {
		flag.Parse()
	}
}
