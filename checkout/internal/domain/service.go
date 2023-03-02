package domain

import "context"

type Service struct {
	lomsServiceClient    LOMSServiceClient
	productServiceClient ProductServiceClient
}

func New(
	lomsServiceClient LOMSServiceClient,
	productServiceClient ProductServiceClient,
) *Service {
	return &Service{
		lomsServiceClient:    lomsServiceClient,
		productServiceClient: productServiceClient,
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

type LOMSServiceClient interface {
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (OrderID, error)
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku uint32) (Product, error)
}
