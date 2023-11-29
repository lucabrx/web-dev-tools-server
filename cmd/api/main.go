package main

import (
	"context"
	"database/sql"

	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wdt/config"
)

var URL = "http://localhost:8080"

type application struct {
	logger *zerolog.Logger
	wg     sync.WaitGroup
	config config.AppConfig
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	db, err := openDB(cfg.DBUrl, 25, 25, 15*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Logger.Info().Msg("Connected to database")

	app := application{
		logger: &log.Logger,
		config: cfg,
	}

	err = app.serve()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")

	}
}

func openDB(url string, maxOpenCon, maxIdleCon int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenCon)
	db.SetMaxIdleConns(maxIdleCon)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}