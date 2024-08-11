package logistics_engine

import (
	"context"
	"errors"
	logistics_v1 "github.com/ivanbulyk/logistics_engine_api/internal/generated/logistics/api/v1"
	"github.com/ivanbulyk/logistics_engine_api/internal/logging"
	"github.com/ivanbulyk/logistics_engine_api/internal/model"
	"github.com/ivanbulyk/logistics_engine_api/internal/repository"
	"log/slog"
	"strconv"
)

type LogisticsEngine struct {
	log          *slog.Logger
	dlvUnitSaver DeliveryUnitSaver
	rptProvider  ReportProvider
}

func NewLogisticsEngine(log *slog.Logger, dlvUnitSaver DeliveryUnitSaver, rptProvider ReportProvider) *LogisticsEngine {
	return &LogisticsEngine{
		log:          log,
		dlvUnitSaver: dlvUnitSaver,
		rptProvider:  rptProvider,
	}
}

type DeliveryUnitSaver interface {
	Create(_ context.Context, report model.MetricsReport) (model.MetricsReport, error)
	Update(_ context.Context, report model.MetricsReport) error
	GetByID(_ context.Context, id int64) (model.MetricsReport, error)
}

type ReportProvider interface {
	GetAll(_ context.Context) ([]model.MetricsReport, error)
}

func (l *LogisticsEngine) MoveUnit(ctx context.Context, in *logistics_v1.MoveUnitRequest) (*logistics_v1.DefaultResponse, error) {
	const opLabel = "LogisticsEngine.MoveUnit"

	log := l.log.With(
		slog.String("opLabel", opLabel),
		slog.String("CargoUnitId", strconv.FormatInt(in.GetCargoUnitId(), 10)),
	)

	report := model.MetricsReport{
		ID: in.GetCargoUnitId(),
		MoveUnit: model.MoveUnit{
			CargoUnitId: in.GetCargoUnitId(),
			Location: []model.Location{
				{
					Latitude:  in.GetLocation().GetLatitude(),
					Longitude: in.GetLocation().GetLongitude(),
				},
			},
		},
		UnitReachedWarehouse: model.UnitReachedWarehouse{
			Location: model.Location{
				Latitude:  0,
				Longitude: 0,
			},
			Announcement: model.WarehouseAnnouncement{
				CargoUnitId: 0,
				WarehouseId: 0,
				Message:     "",
			},
		},
	}

	log.Info("attempting to create and save metrics report with move unit data")

	_, err := l.dlvUnitSaver.Create(ctx, report)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			l.log.Warn("metrics report with the id already exists. updating move unit data", logging.Err(err))
			report, _ = l.dlvUnitSaver.GetByID(ctx, in.GetCargoUnitId())
			location := model.Location{
				Latitude:  in.GetLocation().GetLatitude(),
				Longitude: in.GetLocation().GetLongitude(),
			}
			report.MoveUnit.Location = append(report.MoveUnit.Location, location)
			_ = l.dlvUnitSaver.Update(ctx, report)

		}

	}

	return &logistics_v1.DefaultResponse{}, nil
}

func (l *LogisticsEngine) UnitReachedWarehouse(ctx context.Context, in *logistics_v1.UnitReachedWarehouseRequest) (*logistics_v1.DefaultResponse, error) {
	const opLabel = "LogisticsEngine.UnitReachedWarehouse"

	log := l.log.With(
		slog.String("opLabel", opLabel),
		slog.String("CargoUnitId", strconv.FormatInt(in.GetAnnouncement().GetCargoUnitId(), 10)),
	)
	report := model.MetricsReport{
		ID: in.GetAnnouncement().GetCargoUnitId(),
		MoveUnit: model.MoveUnit{
			CargoUnitId: 0,
			Location:    []model.Location{},
		},
		UnitReachedWarehouse: model.UnitReachedWarehouse{
			Location: model.Location{
				Latitude:  in.GetLocation().GetLatitude(),
				Longitude: in.GetLocation().GetLongitude(),
			},
			Announcement: model.WarehouseAnnouncement{
				CargoUnitId: in.GetAnnouncement().GetCargoUnitId(),
				WarehouseId: in.GetAnnouncement().GetWarehouseId(),
				Message:     in.GetAnnouncement().GetMessage(),
			},
		},
	}

	log.Info("attempting to create and save metrics report with unit reached warehouse data")

	_, err := l.dlvUnitSaver.Create(ctx, report)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			l.log.Warn("metrics report with the id already exists. updating unit reached warehouse data", logging.Err(err))
			report, _ = l.dlvUnitSaver.GetByID(ctx, in.GetAnnouncement().GetCargoUnitId())

			ur := model.UnitReachedWarehouse{
				Location: model.Location{
					Latitude:  in.GetLocation().GetLatitude(),
					Longitude: in.GetLocation().GetLongitude(),
				},
				Announcement: model.WarehouseAnnouncement{
					CargoUnitId: in.GetAnnouncement().GetCargoUnitId(),
					WarehouseId: in.GetAnnouncement().GetWarehouseId(),
					Message:     in.GetAnnouncement().GetMessage(),
				},
			}
			report.UnitReachedWarehouse = ur
			_ = l.dlvUnitSaver.Update(ctx, report)

		}

	}

	return &logistics_v1.DefaultResponse{}, nil
}

func (l *LogisticsEngine) MetricsReport(ctx context.Context, in *logistics_v1.DefaultRequest) (*logistics_v1.MetricsReportResponse, error) {
	const opLabel = "LogisticsEngine.MetricsReport"

	log := l.log.With(
		slog.String("opLabel", opLabel),
	)

	log.Info("attempting to get metrics report")

	report, err := l.rptProvider.GetAll(ctx)
	if err != nil {
		log.Error("failed to get metrics report", logging.Err(err))
		return nil, err
	}

	mr := &logistics_v1.MetricsReportResponse{
		DeliveryUnitsNumber:                           int64(len(report)),
		WarehousesReceivedSuppliesList:                warehousesReceivedSuppliesList(report),
		DeliveryUnitsReachedDestination:               deliveryUnitsReachedDestination(report),
		DeliveryUnitsEachWarehouseReceivedTotalNumber: deliveryUnitsEachWarehouseReceivedTotalNumber(report),
	}

	return mr, nil
}

// warehousesReceivedSuppliesList returns a list of warehouses that have received supplies
func warehousesReceivedSuppliesList(report []model.MetricsReport) []int64 {

	data := make(map[int64][]int64)

	list := []int64{}
	for _, r := range report {
		key := r.UnitReachedWarehouse.Announcement.WarehouseId
		value := r.UnitReachedWarehouse.Announcement.CargoUnitId
		if v, exist := data[key]; exist {
			data[key] = append(v, value)
		} else {
			data[key] = []int64{value}
		}

	}
	for key := range data {
		list = append(list, key)
	}
	return list
}

// deliveryUnitsReachedDestination returns a list units that have reached their destination
func deliveryUnitsReachedDestination(report []model.MetricsReport) []int64 {
	list := []int64{}
	for _, r := range report {
		list = append(list, r.UnitReachedWarehouse.Announcement.CargoUnitId)
	}
	return list
}

// deliveryUnitsEachWarehouseReceivedTotalNumber
func deliveryUnitsEachWarehouseReceivedTotalNumber(report []model.MetricsReport) []*logistics_v1.DeliveryUnitsWarehouseReceivedTotalNumber {
	data := make(map[int64][]int64)

	list := []*logistics_v1.DeliveryUnitsWarehouseReceivedTotalNumber{}
	for _, r := range report {
		key := r.UnitReachedWarehouse.Announcement.WarehouseId
		value := r.UnitReachedWarehouse.Announcement.CargoUnitId
		if v, exist := data[key]; exist {
			data[key] = append(v, value)
		} else {
			data[key] = []int64{value}
		}
	}
	for k, v := range data {
		unit := &logistics_v1.DeliveryUnitsWarehouseReceivedTotalNumber{}
		unit.WarehouseId = k
		unit.DeliveryUnitsNumber = int64(len(v))
		list = append(list, unit)
	}
	return list
}
