package v1

import (
	"context"
	"route256/loms/internal/models"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) ListOrder(ctx context.Context, reqPayload *apiSchema.OrderID) (*apiSchema.ListOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, models.OrderID(reqPayload.OrderID))
	if err != nil {
		return nil, err
	}

	resPayload := apiSchema.ListOrderResponse{
		User:   int64(order.User),
		Status: string(order.Status),
		Items:  make([]*apiSchema.OrderItem, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		resPayload.Items = append(resPayload.Items, &apiSchema.OrderItem{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		})
	}
	return &resPayload, nil
}
