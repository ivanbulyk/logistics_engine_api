package app

import (
	"context"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/ivanbulyk/logistics_engine_api/internal/config"
	"github.com/ivanbulyk/logistics_engine_api/internal/grpcapp"
	"github.com/ivanbulyk/logistics_engine_api/internal/httpapp"
	"github.com/ivanbulyk/logistics_engine_api/internal/logging"
	"github.com/ivanbulyk/logistics_engine_api/internal/repository/memory"
	"github.com/ivanbulyk/logistics_engine_api/internal/services/logistics_engine"
	"golang.org/x/sync/errgroup"
	stndlog "log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type App struct {
	GRPCApp *grpcapp.App
	HTTPApp *httpapp.App
}

// New returns an App instance.
func New(log *slog.Logger, grpcPort int, httpAddr string, srvMetrics *grpcprom.ServerMetrics, reg *prometheus.Registry) *App {

	repository := memory.New()

	logisticsEngineService := logistics_engine.NewLogisticsEngine(log, repository, repository)
	grpcApp := grpcapp.New(grpcPort, log, logisticsEngineService, srvMetrics)
	httpApp := httpapp.New(httpAddr, log, reg)
	return &App{
		GRPCApp: grpcApp,
		HTTPApp: httpApp,
	}
}

// MustRun is wrapper around run() and it panics if any error occurs.
func MustRun() {
	if err := run(); err != nil {
		panic(err)
	}
}

// Run runs App instance.
func run() error {

	cfg := &config.ServerAppConfig{}
	cfg.LoadFromEnv()
	log := logging.SetupLogger(cfg.LogLevel)
	gRPCPort, _ := strconv.Atoi(cfg.Port)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// handle CTRL+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(quit)

	g, ctx := errgroup.WithContext(ctx)

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{1}), // 0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	application := New(log, gRPCPort, config.METRICS_SERVER, srvMetrics, reg)

	g.Go(func() error {
		application.GRPCApp.MustRun()
		return nil
	})

	g.Go(func() error {
		application.HTTPApp.MustRun()

		return nil
	})

	// handle termination
	select {
	case <-quit:
		break
	case <-ctx.Done():
		break
	}

	// gracefully shutdown servers
	cancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer timeoutCancel()

	log.Info("shutting down servers, please wait...")

	application.GRPCApp.Stop()
	application.HTTPApp.Stop(timeoutCtx)

	// wait for shutdown
	if err := g.Wait(); err != nil {
		stndlog.Fatal(err)
	}

	return nil
}
