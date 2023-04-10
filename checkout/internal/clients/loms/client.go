package loms

import (
	"route256/checkout/internal/domain"
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	lomsServiceClient lomsServiceAPI.LomsV1Client
}

func New(address string) (domain.LOMSServiceClient, error) {
	cc, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		lomsServiceClient: lomsServiceAPI.NewLomsV1Client(cc),
	}, nil
}
