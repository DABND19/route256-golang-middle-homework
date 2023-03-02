package product

import (
	"route256/checkout/internal/domain"
	productServiceAPI "route256/product-service/pkg/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	productServiceClient productServiceAPI.ProductServiceClient
	token                string
}

func New(address string, token string) (domain.ProductServiceClient, error) {
	cc, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		productServiceClient: productServiceAPI.NewProductServiceClient(cc),
		token:                token,
	}, nil
}
