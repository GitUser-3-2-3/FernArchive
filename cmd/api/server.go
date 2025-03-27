package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (bknd *backend) serve() error {
	srvr := &http.Server{
		ErrorLog:     slog.NewLogLogger(bknd.logger.Handler(), slog.LevelError),
		Addr:         fmt.Sprintf(":%d", bknd.config.port),
		Handler:      bknd.routes(),
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	bknd.logger.Info("API server started", "addrs", srvr.Addr, "env", bknd.config.env)
	return srvr.ListenAndServe()
}
