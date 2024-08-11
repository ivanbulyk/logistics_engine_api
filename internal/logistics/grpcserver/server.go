package grpcserver

import (
	"context"
	logistics_v1 "github.com/ivanbulyk/logistics_engine_api/internal/generated/logistics/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogisticsEngine interface {
	MoveUnit(ctx context.Context, in *logistics_v1.MoveUnitRequest) (*logistics_v1.DefaultResponse, error)
	UnitReachedWarehouse(ctx context.Context, in *logistics_v1.UnitReachedWarehouseRequest) (*logistics_v1.DefaultResponse, error)
	MetricsReport(ctx context.Context, in *logistics_v1.DefaultRequest) (*logistics_v1.MetricsReportResponse, error)
}

type server struct {
	logistics_v1.UnimplementedLogisticsEngineAPIServer
	logisticsEngine LogisticsEngine
}

func Register(gRPC *grpc.Server, logisticsEngine LogisticsEngine) {

	logistics_v1.RegisterLogisticsEngineAPIServer(gRPC, &server{logisticsEngine: logisticsEngine})
}

func (s *server) MoveUnit(ctx context.Context, in *logistics_v1.MoveUnitRequest) (*logistics_v1.DefaultResponse, error) {
	defaultResponse, err := s.logisticsEngine.MoveUnit(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to process incoming request")
	}

	return defaultResponse, nil
}
func (s *server) UnitReachedWarehouse(ctx context.Context, in *logistics_v1.UnitReachedWarehouseRequest) (*logistics_v1.DefaultResponse, error) {
	defaultResponse, err := s.logisticsEngine.UnitReachedWarehouse(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to process incoming request")
	}

	return defaultResponse, nil
}

func (s *server) MetricsReport(ctx context.Context, in *logistics_v1.DefaultRequest) (*logistics_v1.MetricsReportResponse, error) {
	metricsReportResponse, err := s.logisticsEngine.MetricsReport(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to process response with metrics report")
	}

	return metricsReportResponse, nil
}
