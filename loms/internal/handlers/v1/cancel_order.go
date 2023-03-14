package v1

import (
	"context"
	"route256/loms/internal/models"
	apiSchema "route256/loms/pkg/lomsv1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CancelOrder(ctx context.Context, reqPayload *apiSchema.OrderID) (*emptypb.Empty, error) {
	err := s.service.CancelOrder(ctx, models.OrderID(reqPayload.OrderID))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
