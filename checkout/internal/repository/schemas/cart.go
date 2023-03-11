package schemas

import "time"

type CartItem struct {
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	SKU       int64     `db:"sku"`
	Count     int32     `db:"count"`
}
