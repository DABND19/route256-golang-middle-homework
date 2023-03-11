package product

import (
	"context"
	"route256/checkout/internal/models"
	productServiceAPI "route256/product-service/pkg/product"
)

func (c *Client) GetProduct(ctx context.Context, sku models.SKU) (*models.Product, error) {
	reqPayload := &productServiceAPI.GetProductRequest{
		Token: c.token,
		Sku:   uint32(sku),
	}
	resPayload, err := c.productServiceClient.GetProduct(ctx, reqPayload)
	return &models.Product{
		Name:  resPayload.GetName(),
		Price: resPayload.GetPrice(),
	}, err
}
