package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/krtech-it/metricagent/internal/config"
	config_db "github.com/krtech-it/metricagent/internal/config/db"
	delivery "github.com/krtech-it/metricagent/internal/delivery/http"
	"github.com/krtech-it/metricagent/internal/logger"
	"log"
	"strconv"
)

func main() {
	config.ParseFlags()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}
	db, err := config_db.NewDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Info(err.Error())
	} else {
		defer db.Close()
	}

	router := delivery.NewRouter(logger.Log, cfg, db)
	logger.Log.Info("Listening on port " + strconv.Itoa(cfg.Port))
	err = router.Run(cfg.Host + ":" + strconv.Itoa(cfg.Port))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
