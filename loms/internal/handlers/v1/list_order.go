package v1

import (
	"context"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) ListOrder(ctx context.Context, reqPayload *apiSchema.OrderID) (*apiSchema.ListOrderResponse, error) {
	return &apiSchema.ListOrderResponse{
		User:   123,
		Status: "new",
		Items: []*apiSchema.OrderItem{
			{Sku: 1, Count: 1},
		},
	}, nil
}
