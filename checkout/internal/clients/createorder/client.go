package createorder

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/serviceclient"
)

type Client struct {
	serviceClient *serviceclient.ServiceClient
	endpointPath  string
}

func New(serviceClient *serviceclient.ServiceClient, endpointPath string) *Client {
	return &Client{
		serviceClient: serviceClient,
		endpointPath:  endpointPath,
	}
}

type OrderItemPayload struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type RequestPayload struct {
	Items []OrderItemPayload `json:"items"`
	User  int64              `json:"user"`
}

type ResponsePayload struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(
	ctx context.Context,
	user int64,
	userOrder []domain.OrderItem,
) (domain.OrderID, error) {
	reqPayload := RequestPayload{
		User:  user,
		Items: make([]OrderItemPayload, 0, len(userOrder)),
	}
	for _, item := range userOrder {
		reqPayload.Items = append(reqPayload.Items, OrderItemPayload{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}
	resPayload := ResponsePayload{}
	err := serviceclient.MakeRequest(ctx, c.serviceClient, c.endpointPath, reqPayload, &resPayload)
	return domain.OrderID(resPayload.OrderID), err
}
