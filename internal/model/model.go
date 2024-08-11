package model

type MoveUnit struct {
	CargoUnitId int64      `json:"cargo_unit_id"`
	Location    []Location `json:"location"`
}

// UnitReachedWarehouse contains WarehouseAnnouncement with Location
type UnitReachedWarehouse struct {
	Location     Location              `json:"location"`
	Announcement WarehouseAnnouncement `json:"announcement"`
}

type WarehouseAnnouncement struct {
	// cargo_unit_id is unique id
	CargoUnitId int64 `json:"cargo_unit_id"`
	// warehouse_id is unique id
	WarehouseId int64 `json:"warehouse_id"`
	// the message contains information about the announcement
	Message string `json:"message"`
}

// Location where entity now located in X,Y Axis
type Location struct {
	Latitude  uint32 `json:"latitude"`
	Longitude uint32 `json:"longitude"`
}

type MetricsReport struct {
	ID                   int64                `json:"report_id"`
	MoveUnit             MoveUnit             `json:"move_unit"`
	UnitReachedWarehouse UnitReachedWarehouse `json:"unit_reached_warehouse"`
}
type DeliveryUnitsWarehouseReceivedTotalNumber struct {
	WarehouseId         int64 `json:"warehouse_id"`
	DeliveryUnitsNumber int64 `json:"delivery_units_number"`
}

type Report struct {
	DeliveryUnitsTotalNumber                      int64                                       `json:"delivery_units_total_number"`
	WarehousesReceivedSuppliesList                []int64                                     `json:"warehouses_received_supplies_list"`
	DeliveryUnitsReachedDestination               []int64                                     `json:"delivery_units_reached_destination"`
	DeliveryUnitsEachWarehouseReceivedTotalNumber []DeliveryUnitsWarehouseReceivedTotalNumber `json:"delivery_units_each_warehouse_received_total_number"`
}
