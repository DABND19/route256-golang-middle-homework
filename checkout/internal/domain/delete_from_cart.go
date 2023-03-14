package domain

import (
	"context"
	"route256/checkout/internal/models"
)

func (s *Service) DeleteFromCart(ctx context.Context, user models.User, sku models.SKU, count uint16) error {
	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		updatedCount, err := s.cartsRepository.UpdateCartItemProductsCount(ctx, user, sku, -int32(count))
		if err != nil {
			return err
		}
		if updatedCount < 0 {
			return InvalidProductsCount
		}
		if updatedCount > 0 {
			return nil
		}

		err = s.cartsRepository.DeleteCartItem(ctx, user, sku)
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
