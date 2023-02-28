package v1

import (
	"context"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) DeleteFromCart(ctx context.Context, reqPayload *apiSchema.EditCartRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
