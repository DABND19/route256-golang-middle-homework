package loms

import (
	"context"
	"route256/checkout/internal/domain"
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if status.Code(err) == codes.FailedPrecondition {
			return domain.OrderID(0), errors.Wrap(domain.OrderCreationError, err.Error())
		}
		return domain.OrderID(0), err
	}
	return domain.OrderID(resPayload.GetOrderID()), nil
}
