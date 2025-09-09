package main

import (
	"flag"
	"log"
	"os"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/server/http"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	var store storage.Storage
	if cfg.Storage.Type == "sql" {
		store, err = storage.NewSQLStorage(cfg.Storage.DSN)
		if err != nil {
			log.Fatalf("Failed to create SQL storage: %v", err)
		}
	} else {
		store = storage.NewMemoryStorage()
	}

	calendarApp := app.New(logg, store)
	httpServer := http.NewServer(cfg.Server.Host, cfg.Server.Port, calendarApp, logg)

	logg.Info("Starting calendar service...")
	if err := httpServer.Start(); err != nil {
		logg.Error("Failed to start server: " + err.Error())
		os.Exit(1)
	}
}
