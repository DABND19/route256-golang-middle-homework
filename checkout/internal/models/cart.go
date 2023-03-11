package models

type User int64

type WarehouseID int64

type ProductsCount uint16

type Stock struct {
	WarehouseID WarehouseID
	Count       uint64
}

type OrderID int64

type SKU uint32

type CartItem struct {
	SKU   SKU
	Count ProductsCount
}

type Product struct {
	Name  string
	Price uint32
}

type CartProduct struct {
	CartItem
	Product
}
