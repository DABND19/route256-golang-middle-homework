package domain

import (
	"context"
	"route256/checkout/internal/models"

	"github.com/pkg/errors"
)

func (s *Service) ListCart(ctx context.Context, user models.User) ([]models.CartItem, error) {
	userOrder, err := s.GetUserOrder(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query user order")
	}

	cartItems := make([]models.CartItem, 0, len(userOrder))
	for _, orderItem := range userOrder {
		product, err := s.productServiceClient.GetProduct(ctx, orderItem.SKU)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to request product")
		}
		cartItems = append(cartItems, models.CartItem{
			SKU:   orderItem.SKU,
			Count: orderItem.Count,
			Name:  product.Name,
			Price: product.Price,
		})
	}
	return cartItems, nil
}

func (s *Service) CalculateTotalPrice(cart []models.CartItem) (total uint32) {
	for _, item := range cart {
		total += item.Price * uint32(item.Count)
	}
	return
}
