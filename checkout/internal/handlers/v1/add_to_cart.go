package v1

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AddToCart(ctx context.Context, reqPayload *apiSchema.EditCartRequest) (*emptypb.Empty, error) {
	err := s.service.AddToCart(ctx, reqPayload.GetUser(), reqPayload.GetSku(), uint16(reqPayload.GetCount()))
	if errors.Is(err, domain.InsufficientStocksError) {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &emptypb.Empty{}, err
}
