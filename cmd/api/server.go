package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		bknd.logger.Info("shutting down server", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srvr.Shutdown(ctx)
	}()
	bknd.logger.Info("server started", "addrs", srvr.Addr, "env", bknd.config.env)
	err := srvr.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}
	bknd.logger.Info("server stopped", "addrs", srvr.Addr, "env", bknd.config.env)
	return nil
}
