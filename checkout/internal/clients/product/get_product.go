package product

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	productServiceAPI "route256/product-service/pkg/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetProduct(ctx context.Context, sku models.SKU) (*models.Product, error) {
	err := c.limiter.Wait(ctx)
	if err != nil {
		return nil, domain.ProductServiceRateLimitError
	}

	reqPayload := &productServiceAPI.GetProductRequest{
		Token: c.token,
		Sku:   uint32(sku),
	}
	resPayload, err := c.productServiceClient.GetProduct(ctx, reqPayload)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ProductNotFound
		}
		return nil, err
	}

	return &models.Product{
		Name:  resPayload.GetName(),
		Price: resPayload.GetPrice(),
	}, nil
}
