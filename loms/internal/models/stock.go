package models

type SKU uint32

type WarehouseID int64

type Stock struct {
	WarehouseID WarehouseID
	Count       uint64
}

type ItemBooking struct {
	WarehouseID WarehouseID
	Count       uint16
}
