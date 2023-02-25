package domain

import (
	"context"
	"errors"
)

var (
	InsufficientStocksError = errors.New("Insufficient stocks")
)

func (s *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	stocks, err := s.stocksChecker.Stocks(ctx, sku)
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
