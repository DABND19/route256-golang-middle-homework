package loms

import (
	"context"
	"route256/checkout/internal/domain"
	lomsServiceAPI "route256/loms/pkg/lomsv1"
)

func (c *Client) CreateOrder(
	ctx context.Context,
	user int64,
	userOrder []domain.OrderItem,
) (domain.OrderID, error) {
	reqPayload := &lomsServiceAPI.CreateOrderRequest{
		User:  user,
		Items: make([]*lomsServiceAPI.OrderItem, 0, len(userOrder)),
	}
	for _, item := range userOrder {
		reqPayload.Items = append(reqPayload.Items, &lomsServiceAPI.OrderItem{
			Sku:   item.SKU,
			Count: uint32(item.Count),
		})
	}
	resPayload, err := c.lomsServiceClient.CreateOrder(ctx, reqPayload)
	if err != nil {
		return domain.OrderID(0), err
	}
	return domain.OrderID(resPayload.GetOrderID()), nil
}
