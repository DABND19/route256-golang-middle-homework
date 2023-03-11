package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/models"
)

var (
	InsufficientStocksError = errors.New("Insufficient stocks")
)

func (s *Service) checkStocks(
	ctx context.Context,
	user models.User,
	sku models.SKU,
	count models.ProductsCount,
) error {
	stocks, err := s.lomsServiceClient.Stocks(ctx, sku)
	if err != nil {
		return errors.New("Failed to check stocks")
	}
	var total uint64 = 0
	for _, stock := range stocks {
		total += stock.Count
		if total >= uint64(count) {
			return nil
		}
	}
	return InsufficientStocksError
}

func (s *Service) AddToCart(
	ctx context.Context,
	user models.User,
	sku models.SKU,
	count models.ProductsCount,
) error {
	if err := s.checkStocks(ctx, user, sku, count); err != nil {
		return err
	}

	err := s.RunRepeatableRead(ctx, func(ctx context.Context) error {
		err := s.cartsRepository.CreateCartItem(ctx, user, sku, count)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
