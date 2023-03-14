package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) CreateOrder(
	ctx context.Context,
	user models.User,
	userOrder []models.CartItem,
) (*models.OrderID, error) {
	reqPayload := &lomsServiceAPI.CreateOrderRequest{
		User:  int64(user),
		Items: make([]*lomsServiceAPI.OrderItem, 0, len(userOrder)),
	}
	for _, item := range userOrder {
		reqPayload.Items = append(reqPayload.Items, &lomsServiceAPI.OrderItem{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		})
	}

	resPayload, err := c.lomsServiceClient.CreateOrder(ctx, reqPayload)
	if err != nil {
		if status.Code(err) == codes.FailedPrecondition {
			return nil, errors.Wrap(domain.OrderCreationError, err.Error())
		}
		return nil, err
	}

	orderID := models.OrderID(resPayload.OrderID)
	return &orderID, nil
}
