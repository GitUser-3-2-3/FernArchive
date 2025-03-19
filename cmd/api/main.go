package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = 1.0

type config struct {
	port int
	env  string
}

type backend struct {
	config config
	logger *slog.Logger
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.env, "env", "development", "Environment")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	bknd := &backend{
		logger: logger,
		config: cfg,
	}
	router := http.NewServeMux()
	router.HandleFunc("/v1/healthcheck", bknd.healthcheckHandler)

	srvr := &http.Server{
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("API server started", "addrs", srvr.Addr, "env", cfg.env)
	err := srvr.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
