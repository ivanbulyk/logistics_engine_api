package httpapp

import (
	"context"
	"fmt"
	"github.com/ivanbulyk/logistics_engine_api/internal/logistics/httpserver"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	httpAddr   string
}

// New creates new http server app.
func New(httpAddr string, log *slog.Logger, reg *prometheus.Registry) *App {

	httpServer := httpserver.NewMetricsServer(httpAddr, reg)

	return &App{
		log:        log,
		httpServer: httpServer,
		httpAddr:   httpAddr,
	}
}

// MustRun runs HTTP server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs HTTP server.
func (a *App) Run() error {
	const opLabel = "httpapp.Run"

	a.log.Info("metrics server listening at %s\n", slog.String("addr", a.httpAddr))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Info("failed to serve metrics: %v\n", err)
		return fmt.Errorf("%s: %w", opLabel, err)
	}

	return nil
}

// Stop stops HTTP server.
func (a *App) Stop(timeoutCtx context.Context) {
	const opLabel = "httpapp.Stop"

	a.log.With(slog.String("opLabel", opLabel)).
		Info("metrics server shutdown")

	a.httpServer.Shutdown(timeoutCtx)
}
