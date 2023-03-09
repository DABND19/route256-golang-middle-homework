package domain

import (
	"context"
	"errors"
	"route256/loms/internal/models"
)

var (
	OrderNotFoundError      = errors.New("Order not found")
	InsufficientStocksError = errors.New("Insufficient stocks")
	StockNotFoundError      = errors.New("Stock not found")
)

type Service struct {
	TransactionRunner
	OrdersRespository
	StocksRespository
}

func New(
	tr TransactionRunner,
	ordersRepo OrdersRespository,
	stocksRepo StocksRespository,
) *Service {
	return &Service{
		tr, ordersRepo, stocksRepo,
	}
}

type TransactionRunner interface {
	RunReadCommited(ctx context.Context, txFn func(ctx context.Context) error) error
	RunRepeatableRead(ctx context.Context, txFn func(ctx context.Context) error) error
	RunSerializable(ctx context.Context, txFn func(ctx context.Context) error) error
	RunInSavepoint(ctx context.Context, txFn func(ctx context.Context) error) error
}

type OrdersRespository interface {
	CreateOrder(ctx context.Context, user models.User, items []models.OrderItem) (*models.OrderID, error)
	GetOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error)
	ChangeOrderStatus(ctx context.Context, orderID models.OrderID, status models.OrderStatus) error
}

type StocksRespository interface {
	GetStocks(ctx context.Context, sku models.SKU) ([]models.Stock, error)
	UpdateStockItemsCount(ctx context.Context, warehouseID models.WarehouseID, sku models.SKU, diff int64) error

	GetItemBookings(ctx context.Context, orderID models.OrderID, sku models.SKU) ([]models.ItemBooking, error)
	CreateItemBooking(ctx context.Context, orderID models.OrderID, warehouseID models.WarehouseID, sku models.SKU, count uint16) error
	DeleteItemBooking(ctx context.Context, orderID models.OrderID, warehouseID models.WarehouseID, sku models.SKU) error
}
