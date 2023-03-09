package schemas

import "time"

type ItemStock struct {
	WarehouseID int64  `db:"warehouse_id"`
	SKU         int64  `db:"sku"`
	Count       uint64 `db:"count"`
}

type ItemBooking struct {
	CreatedAt   time.Time `db:"created_at"`
	OrderID     int64     `db:"order_id"`
	WarehouseID int64     `db:"warehouse_id"`
	SKU         int64     `db:"sku"`
	Count       uint16    `db:"count"`
}
