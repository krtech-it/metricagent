package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/krtech-it/metricagent/internal/config"
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
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := delivery.NewRouter(logger.Log, cfg, db)
	logger.Log.Info("Listening on port " + strconv.Itoa(cfg.Port))
	err = router.Run(cfg.Host + ":" + strconv.Itoa(cfg.Port))
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}
