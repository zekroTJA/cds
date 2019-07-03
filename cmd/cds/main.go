package main

import (
	"flag"
	"os"

	"github.com/zekroTJA/cds/internal/config"
	"github.com/zekroTJA/cds/internal/database"
	"github.com/zekroTJA/cds/internal/logger"
	"github.com/zekroTJA/cds/internal/static"
	"github.com/zekroTJA/cds/internal/webserver"
)

var (
	flagConfig = flag.String("c", "config.yml", "config file location")
	flagAddr   = flag.String("addr", "", "expose address (overrides config)")
)

func main() {
	flag.Parse()

	logger.Info("cds version %s", static.AppVersion)

	logger.Info("CONFIG :: initializing...")
	cfg, err := config.Open(*flagConfig)
	if err != nil {
		logger.Error("CONFIG :: failed loading: %s", err.Error())
		os.Exit(1)
	} else if cfg == nil {
		logger.Error("CONFIG :: config file was not found. A defautl config file was created. " +
			"Edit this file and restart after.")
		os.Exit(0)
	}

	if *flagAddr != "" {
		cfg.WebServer.Addr = *flagAddr
	}

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		logger.Error("DATABASE :: failed connection to database: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		logger.Info("DATABASE :: tear down")
		db.Close()
	}()

	logger.Info("WEBSERVER :: webserver started on address %s", cfg.WebServer.Addr)
	err = webserver.NewWebServer(cfg.WebServer, db).
		ListenAndServeBlocking()
	logger.Error("WEBSERVER :: failed starting web server: %s", err.Error())
	defer logger.Info("WEBSERVER :: tear down")
}
