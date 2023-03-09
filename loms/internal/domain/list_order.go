package domain

import (
	"context"
	"route256/loms/internal/models"
)

func (s *Service) ListOrder(ctx context.Context, orderID models.OrderID) (*models.Order, error) {
	var order *models.Order
	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		var err error
		order, err = s.OrdersRespository.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
