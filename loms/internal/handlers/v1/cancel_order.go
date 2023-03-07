package v1

import (
	"context"
	apiSchema "route256/loms/pkg/lomsv1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CancelOrder(ctx context.Context, reqPayload *apiSchema.OrderID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
