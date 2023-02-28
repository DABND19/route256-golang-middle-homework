package loms

import (
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"google.golang.org/grpc"
)

type Client struct {
	lomsServiceClient lomsServiceAPI.LomsV1Client
}

func New(cc *grpc.ClientConn) *Client {
	return &Client{
		lomsServiceClient: lomsServiceAPI.NewLomsV1Client(cc),
	}
}
