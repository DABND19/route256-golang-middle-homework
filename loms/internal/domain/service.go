package domain

import (
	"context"
	"errors"
	"route256/loms/internal/models"
	"time"
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
	NotificationsClient
	OrderStatusChangeRepository
	cancelOrderScheduler              Scheduler
	unpaidOrderTtl                    time.Duration
	ordersStatusChangesSumbitInterval time.Duration
}

func New(
	ctx context.Context,
	tr TransactionRunner,
	ordersRepo OrdersRespository,
	stocksRepo StocksRespository,
	unpaidOrderTtl time.Duration,
	cancelOrderScheduler Scheduler,
	orderStatusChangeNotifier NotificationsClient,
	orderStatusChangeRepository OrderStatusChangeRepository,
	ordersStatusChangesSumbitInterval time.Duration,
) *Service {
	s := &Service{
		tr,
		ordersRepo,
		stocksRepo,
		orderStatusChangeNotifier,
		orderStatusChangeRepository,
		cancelOrderScheduler,
		unpaidOrderTtl,
		ordersStatusChangesSumbitInterval,
	}
	s.runOrdersStatusChangesSubmission(ctx)
	return s
}

type TransactionRunner interface {
	RunReadCommited(ctx context.Context, txFn func(ctx context.Context) error) error
	RunRepeatableRead(ctx context.Context, txFn func(ctx context.Context) error) error
	RunSerializable(ctx context.Context, txFn func(ctx context.Context) error) error
	RunInSavepoint(ctx context.Context, txFn func(ctx context.Context) error) error
}

type Scheduler interface {
	Schedule(after time.Time, task func())
}

type OrdersRespository interface {
	CreateOrder(ctx context.Context, user models.User, items []models.OrderItem) (*models.OrderID, error)
	GetOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID models.OrderID, status models.OrderStatus) error
}

type StocksRespository interface {
	GetStocks(ctx context.Context, sku models.SKU) ([]models.Stock, error)
	UpdateStockItemsCount(ctx context.Context, warehouseID models.WarehouseID, sku models.SKU, diff int64) error

	GetItemBookings(ctx context.Context, orderID models.OrderID, sku models.SKU) ([]models.ItemBooking, error)
	CreateItemBooking(ctx context.Context, orderID models.OrderID, warehouseID models.WarehouseID, sku models.SKU, count uint16) error
	DeleteItemBooking(ctx context.Context, orderID models.OrderID, warehouseID models.WarehouseID, sku models.SKU) error
}

type OrderStatusChangeRepository interface {
	GetUnsubmittedChanges(ctx context.Context) ([]models.OrderStatusChange, error)
	MarkChangeAsSubmitted(ctx context.Context, orderStatusChange models.OrderStatusChange) error
	LogOrderStatusChange(ctx context.Context, orderID models.OrderID, status models.OrderStatus) error
}

type NotificationsClient interface {
	NotifyAboutOrderStatusChange(ctx context.Context, change models.OrderStatusChange) error
}
