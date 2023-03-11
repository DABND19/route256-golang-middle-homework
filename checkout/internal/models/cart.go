package models

type User int64

type WarehouseID int64

type Stock struct {
	WarehouseID WarehouseID
	Count       uint64
}

type Product struct {
	Name  string
	Price uint32
}

type OrderID int64

type SKU uint32

type OrderItem struct {
	SKU   SKU
	Count uint16
}

type CartItem struct {
	SKU   SKU
	Count uint16
	Name  string
	Price uint32
}
