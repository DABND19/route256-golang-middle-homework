package domain

import (
	"context"

	"github.com/pkg/errors"
)

type CartItem struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

func (s *Service) ListCart(ctx context.Context, user int64) ([]CartItem, error) {
	userOrder, err := s.GetUserOrder(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query user order")
	}

	cartItems := make([]CartItem, 0, len(userOrder))
	for _, orderItem := range userOrder {
		product, err := s.productGetter.GetProduct(ctx, orderItem.SKU)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to request product")
		}
		cartItems = append(cartItems, CartItem{
			SKU:   orderItem.SKU,
			Count: orderItem.Count,
			Name:  product.Name,
			Price: product.Price,
		})
	}
	return cartItems, nil
}

func (s *Service) CalculateTotalPrice(cart []CartItem) (total uint32) {
	for _, item := range cart {
		total += item.Price * uint32(item.Count)
	}
	return
}
