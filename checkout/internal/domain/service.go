package domain

import "context"

type Service struct {
	stocksChecker StocksChecker
	productGetter ProductGetter
	orderCreator  OrderCreator
}

func New(
	stocksChecker StocksChecker,
	productGetter ProductGetter,
	orderCreator OrderCreator,
) *Service {
	return &Service{
		stocksChecker: stocksChecker,
		productGetter: productGetter,
		orderCreator:  orderCreator,
	}
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Product struct {
	Name  string
	Price uint32
}

type OrderID int64

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type StocksChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductGetter interface {
	GetProduct(ctx context.Context, sku uint32) (Product, error)
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error)
}
