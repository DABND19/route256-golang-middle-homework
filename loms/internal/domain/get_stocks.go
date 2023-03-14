package domain

import (
	"context"
	"route256/loms/internal/models"
)

func (s *Service) GetStocks(ctx context.Context, sku models.SKU) ([]models.Stock, error) {
	var stocks []models.Stock

	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		var err error
		stocks, err = s.StocksRespository.GetStocks(ctx, sku)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
