package domain

import (
	"context"
	"route256/checkout/internal/models"

	"github.com/pkg/errors"
)

var (
	OrderCreationError = errors.New("Failed to create order")
)

func (s *Service) MakePurchase(ctx context.Context, user models.User) (*models.OrderID, error) {
	userOrder, err := s.GetUserOrder(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query user order")
	}

	orderID, err := s.lomsServiceClient.CreateOrder(ctx, user, userOrder)
	if err != nil {
		return nil, err
	}
	return orderID, nil
}
