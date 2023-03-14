package schemas

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
