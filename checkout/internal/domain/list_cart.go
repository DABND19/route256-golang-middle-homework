package domain

import (
	"context"
	"route256/checkout/internal/models"

	"github.com/pkg/errors"
)

var (
	ProductServiceRateLimitError = errors.New("Too many requests to product service")
)

func (s *Service) ListCart(ctx context.Context, user models.User) ([]models.CartProduct, error) {
	var cartItems []models.CartItem
	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		var err error
		cartItems, err = s.cartsRepository.GetCartItems(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query user cart")
	}

	cartProducts := make([]models.CartProduct, 0, len(cartItems))
	for _, item := range cartItems {
		product, err := s.productServiceClient.GetProduct(ctx, item.SKU)
		if err != nil {
			if errors.Is(err, ProductNotFound) {
				continue
			}
			return nil, errors.Wrap(err, "Failed to request product")
		}

		cartProducts = append(cartProducts, models.CartProduct{
			CartItem: item,
			Product:  *product,
		})
	}
	return cartProducts, nil
}

func (s *Service) CalculateTotalPrice(cart []models.CartProduct) (total uint32) {
	for _, item := range cart {
		total += item.Price * uint32(item.Count)
	}
	return
}
