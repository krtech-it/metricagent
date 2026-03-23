package config

import (
	"flag"
)

type SetServer struct {
	addr            string
	storeInterval   int
	fileStoragePath string
	restore         bool
	databaseDSN     string
	hashKey         string
}

var FlagServer = SetServer{}

func ParseFlags() {
	flag.StringVar(&FlagServer.addr, "a", "localhost:8080", "server listen address")
	flag.IntVar(&FlagServer.storeInterval, "i", 300, "interval in seconds")
	flag.StringVar(&FlagServer.fileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&FlagServer.restore, "r", false, "restore storage")
	flag.StringVar(&FlagServer.databaseDSN, "d", "", "database DSN")
	flag.StringVar(&FlagServer.hashKey, "k", "", "hash key")
	if !flag.Parsed() {
		flag.Parse()
	}
}
