package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"FernArchive/internal/data"
	"FernArchive/internal/mailer"
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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		allowedOrigins []string
	}
}

type backend struct {
	logger *slog.Logger
	config config
	models data.Models
	mailer mailer.Mailer
	wtgrp  sync.WaitGroup
}

func main() {
	var cfg config
	runClFlags(&cfg)

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
		models: data.NewModels(db),
		mailer: mailer.NewMailer(cfg.smtp.host, cfg.smtp.port,
			cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	err = bknd.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func runClFlags(cfg *config) {
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev, staging, prod)")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn",
		"postgres://archive:Qwerty1,0*@localhost/archive_db?sslmode=disable", "PostgresSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "DB max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "DB max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "DB max idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Limiter max requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Limiter max burst requests")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiting")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "1702867f97eaf8", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "84bdbfed10e5d6", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender",
		"FernArchive <parthsrivastav.00@gmail.com>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.allowedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()
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
