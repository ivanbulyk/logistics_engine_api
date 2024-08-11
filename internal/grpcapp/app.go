package grpcapp

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/ivanbulyk/logistics_engine_api/internal/logistics/grpcserver"
	"log/slog"
	"net"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app.
func New(port int, log *slog.Logger, logisticsEngine grpcserver.LogisticsEngine, srvMetrics *grpcprom.ServerMetrics) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			srvMetrics.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(grpcserver.InterceptorLogger(log), loggingOpts...),
		),
	}

	gRPCServer := grpc.NewServer(
		opts...,
	)
	grpcserver.Register(gRPCServer, logisticsEngine)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const opLabel = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", opLabel, err)
	}

	a.log.Info("gRPC server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", opLabel, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const opLabel = "grpcapp.Stop"

	a.log.With(slog.String("opLabel", opLabel)).
		Info("gRPC server shutdown")

	a.gRPCServer.GracefulStop()
}
