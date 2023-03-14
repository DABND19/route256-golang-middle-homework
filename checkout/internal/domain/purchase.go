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
	var orderID *models.OrderID
	err := s.RunSerializable(ctx, func(ctx context.Context) error {
		var err error
		cartItems, err := s.cartsRepository.GetCartItems(ctx, user)
		if err != nil {
			return err
		}

		orderID, err = s.lomsServiceClient.CreateOrder(ctx, user, cartItems)
		if err != nil {
			return err
		}

		for _, item := range cartItems {
			err := s.cartsRepository.DeleteCartItem(ctx, user, item.SKU)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return orderID, nil
}
