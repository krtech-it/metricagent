package config

import (
	"flag"
)

type SetServer struct {
	addr            string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

var FlagServer = SetServer{}

func ParseFlags() {
	flag.StringVar(&FlagServer.addr, "a", "localhost:8080", "server listen address")
	flag.IntVar(&FlagServer.storeInterval, "i", 300, "interval in seconds")
	flag.StringVar(&FlagServer.fileStoragePath, "f", "storage.json", "file storage path")
	flag.BoolVar(&FlagServer.restore, "r", false, "restore storage")
	if !flag.Parsed() {
		flag.Parse()
	}
}
