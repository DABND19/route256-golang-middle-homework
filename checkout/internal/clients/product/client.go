package product

import (
	productServiceAPI "route256/product-service/pkg/product"

	"google.golang.org/grpc"
)

type Client struct {
	productServiceClient productServiceAPI.ProductServiceClient
	token                string
}

func New(cc *grpc.ClientConn, token string) *Client {
	return &Client{
		productServiceClient: productServiceAPI.NewProductServiceClient(cc),
		token:                token,
	}
}
