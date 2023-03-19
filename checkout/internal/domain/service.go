package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/models"
)

type Service struct {
	TransactionRunner
	cartsRepository      CartsRepository
	lomsServiceClient    LOMSServiceClient
	productServiceClient ProductServiceClient
}

func New(
	tr TransactionRunner,
	cartsRepo CartsRepository,
	lomsServiceClient LOMSServiceClient,
	productServiceClient ProductServiceClient,
) *Service {
	return &Service{
		tr,
		cartsRepo,
		lomsServiceClient,
		productServiceClient,
	}
}

var (
	CartItemNotFoundError = errors.New("Cart item not found")
	InvalidProductsCount  = errors.New("Invalid number of products")
	ProductNotFound       = errors.New("Product not found")
)

type LOMSServiceClient interface {
	CreateOrder(ctx context.Context, user models.User, items []models.CartItem) (*models.OrderID, error)
	Stocks(ctx context.Context, sku models.SKU) ([]models.Stock, error)
}

type ProductServiceClient interface {
	GetProducts(ctx context.Context, skus []models.SKU) (map[models.SKU]*models.Product, error)
}

type TransactionRunner interface {
	RunReadCommited(ctx context.Context, txFn func(ctx context.Context) error) error
	RunRepeatableRead(ctx context.Context, txFn func(ctx context.Context) error) error
	RunSerializable(ctx context.Context, txFn func(ctx context.Context) error) error
	RunInSavepoint(ctx context.Context, txFn func(ctx context.Context) error) error
}

type CartsRepository interface {
	GetCartItems(ctx context.Context, user models.User) ([]models.CartItem, error)
	CreateCartItem(ctx context.Context, user models.User, sku models.SKU, count models.ProductsCount) error
	UpdateCartItemProductsCount(ctx context.Context, user models.User, sku models.SKU, diff int32) (int32, error)
	DeleteCartItem(ctx context.Context, user models.User, sku models.SKU) error
}
