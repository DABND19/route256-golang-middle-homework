package v1

import (
	"context"
	"route256/checkout/internal/models"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AddToCart(ctx context.Context, reqPayload *apiSchema.EditCartRequest) (*emptypb.Empty, error) {
	err := s.service.AddToCart(ctx, models.User(reqPayload.User), models.SKU(reqPayload.Sku), uint16(reqPayload.GetCount()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}
