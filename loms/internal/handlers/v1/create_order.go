package v1

import (
	"context"
	"route256/loms/internal/models"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) CreateOrder(ctx context.Context, reqPayload *apiSchema.CreateOrderRequest) (*apiSchema.OrderID, error) {
	orderItems := make([]models.OrderItem, 0, len(reqPayload.Items))
	for _, item := range reqPayload.Items {
		orderItems = append(orderItems, models.OrderItem{
			SKU:   models.SKU(item.Sku),
			Count: uint16(item.Count),
		})
	}

	orderID, err := s.service.CreateOrder(ctx, models.User(reqPayload.User), orderItems)
	if err != nil {
		return nil, err
	}
	return &apiSchema.OrderID{OrderID: int64(*orderID)}, nil
}
