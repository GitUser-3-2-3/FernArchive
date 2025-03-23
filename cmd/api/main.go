package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

import _ "github.com/lib/pq"

const version = "1.0.0"

type config struct {
	env  string
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
		maxIdleConns int
	}
}

type backend struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev, staging, prod)")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("ARCHIVE_DB_DSN"), "PostgresSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "db max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "db max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "db max idle time")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(db)
	logger.Info("Database connection established")
	bknd := &backend{
		logger: logger,
		config: cfg,
	}
	srvr := &http.Server{
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      bknd.routes(),
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("API server started", "addrs", srvr.Addr, "env", cfg.env)
	err = srvr.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	} else {
		return db, nil
	}
}
