package v1

import (
	"context"
	"route256/checkout/internal/models"
	apiSchema "route256/checkout/pkg/checkoutv1"
)

func (s *Service) ListCart(ctx context.Context, reqPayload *apiSchema.User) (*apiSchema.Cart, error) {
	cart, err := s.service.ListCart(ctx, models.User(reqPayload.User))
	if err != nil {
		return nil, err
	}
	totalPrice := s.service.CalculateTotalPrice(cart)

	resPayload := &apiSchema.Cart{
		TotalPrice: totalPrice,
		Items:      make([]*apiSchema.CartItem, 0, len(cart)),
	}
	for _, item := range cart {
		resPayload.Items = append(resPayload.Items, &apiSchema.CartItem{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
	}
	return resPayload, nil
}
