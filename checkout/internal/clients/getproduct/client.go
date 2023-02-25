package getproduct

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/serviceclient"
)

type Client struct {
	serviceClient *serviceclient.ServiceClient
	endpointPath  string
	token         string
}

func New(serviceClient *serviceclient.ServiceClient, endpointPath string, token string) *Client {
	return &Client{
		serviceClient: serviceClient,
		endpointPath:  endpointPath,
		token:         token,
	}
}

type RequestPayload struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type ResponsePayload struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (domain.Product, error) {
	reqPayload := RequestPayload{
		Token: c.token,
		SKU:   sku,
	}
	resPayload := ResponsePayload{}
	err := c.serviceClient.Request(ctx, c.endpointPath, reqPayload, &resPayload)
	return domain.Product{
		Name:  resPayload.Name,
		Price: resPayload.Price,
	}, err
}
