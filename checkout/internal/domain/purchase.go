package domain

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Service) MakePurchase(ctx context.Context, user int64) (OrderID, error) {
	userOrder, err := s.GetUserOrder(ctx, user)
	if err != nil {
		return OrderID(0), errors.Wrap(err, "Failed to query user order")
	}

	orderID, err := s.orderCreator.CreateOrder(ctx, user, userOrder)
	if err != nil {
		return OrderID(0), errors.Wrap(err, "Failed to request order creation")
	}
	return orderID, nil
}
