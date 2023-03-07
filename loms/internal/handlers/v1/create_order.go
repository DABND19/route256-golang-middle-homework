package v1

import (
	"context"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) CreateOrder(ctx context.Context, reqPayload *apiSchema.CreateOrderRequest) (*apiSchema.OrderID, error) {
	return &apiSchema.OrderID{OrderID: 1}, nil
}
