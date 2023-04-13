package product

import (
	"route256/checkout/internal/domain"
	"route256/libs/workerpool"
	productServiceAPI "route256/product-service/pkg/product"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerPool interface {
	Submit(task func())
}

type Client struct {
	productServiceClient productServiceAPI.ProductServiceClient
	token                string
	limiter              *rate.Limiter
	wp                   WorkerPool
}

func New(address string, token string, rateLimit int, workerPool workerpool.WorkerPool) (domain.ProductServiceClient, error) {
	cc, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		productServiceClient: productServiceAPI.NewProductServiceClient(cc),
		token:                token,
		limiter:              rate.NewLimiter(rate.Every(1*time.Second), rateLimit),
		wp:                   workerPool,
	}, nil
}
