package product

import (
	"context"
	"route256/checkout/internal/domain"
	productServiceAPI "route256/product-service/pkg/product"
)

func (c *Client) GetProduct(ctx context.Context, sku uint32) (domain.Product, error) {
	reqPayload := &productServiceAPI.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	}
	resPayload, err := c.productServiceClient.GetProduct(ctx, reqPayload)
	return domain.Product{
		Name:  resPayload.GetName(),
		Price: resPayload.GetPrice(),
	}, err
}
