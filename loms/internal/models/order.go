package models

import "time"

type User int64

type OrderID int64

type OrderItem struct {
	SKU   SKU
	Count uint16
}

type OrderStatus string

const (
	OrderStatusNew             = OrderStatus("new")
	OrderStatusAwaitingPayment = OrderStatus("awaiting payment")
	OrderStatusFailed          = OrderStatus("failed")
	OrderStatusPayed           = OrderStatus("payed")
	OrderStatusCancelled       = OrderStatus("cancelled")
)

type Order struct {
	Status OrderStatus
	User   User
	Items  []OrderItem
}

type OrderStatusChange struct {
	CreatedAt time.Time
	OrderID   OrderID
	Status    OrderStatus
}
