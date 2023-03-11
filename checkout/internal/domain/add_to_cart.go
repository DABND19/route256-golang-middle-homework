package domain

import (
	"context"
	"errors"
	"route256/checkout/internal/models"
)

var (
	InsufficientStocksError = errors.New("Insufficient stocks")
)

func (s *Service) AddToCart(ctx context.Context, user models.User, sku models.SKU, count uint16) error {
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
