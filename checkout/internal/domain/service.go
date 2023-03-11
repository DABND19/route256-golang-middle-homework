package domain

import (
	"context"
	"route256/checkout/internal/models"
)

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

type LOMSServiceClient interface {
	CreateOrder(ctx context.Context, user models.User, items []models.OrderItem) (*models.OrderID, error)
	Stocks(ctx context.Context, sku models.SKU) ([]models.Stock, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku models.SKU) (*models.Product, error)
}
