package product

import (
	"route256/checkout/internal/domain"
	productServiceAPI "route256/product-service/pkg/product"
	"time"

	"golang.org/x/time/rate"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	productServiceClient productServiceAPI.ProductServiceClient
	token                string
	limiter              *rate.Limiter
}

func New(address string, token string, rateLimit int) (domain.ProductServiceClient, error) {
	cc, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		productServiceClient: productServiceAPI.NewProductServiceClient(cc),
		token:                token,
		limiter:              rate.NewLimiter(rate.Every(1*time.Second), rateLimit),
	}, nil
}
