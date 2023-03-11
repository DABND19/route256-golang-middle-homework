package v1

import (
	"context"
	"route256/checkout/internal/models"
	"route256/checkout/internal/validators"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) DeleteFromCart(ctx context.Context, reqPayload *apiSchema.EditCartRequest) (*emptypb.Empty, error) {
	if err := validators.ValidateProductsCount(reqPayload); err != nil {
		return nil, err
	}

	err := s.service.DeleteFromCart(ctx, models.User(reqPayload.User), models.SKU(reqPayload.Sku), uint16(reqPayload.Count))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
