package schemas

import "time"

const (
	OrderStatusNew             = "NEW"
	OrderStatusAwaitingPayment = "AWAITING_PAYMENT"
	OrderStatusFailed          = "FAILED"
	OrderStatusPayed           = "PAYED"
	OrderStatusCancelled       = "CANCELLED"
)

type Order struct {
	OrderID int64  `db:"order_id"`
	UserID  int64  `db:"user_id"`
	Status  string `db:"status"`
}

type OrderItem struct {
	OrderID int64  `db:"order_id"`
	SKU     int64  `db:"sku"`
	Count   uint16 `db:"count"`
}

type OrderStatusChange struct {
	CreatedAt   time.Time  `db:"created_at"`
	SubmittedAt *time.Time `db:"submitted_at"`
	OrderID     int64      `db:"order_id"`
	Status      string     `db:"status"`
}
