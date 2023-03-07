package v1

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Purchase(ctx context.Context, reqPayload *apiSchema.User) (*apiSchema.OrderID, error) {
	orderID, err := s.service.MakePurchase(ctx, reqPayload.User)
	if err != nil {
		if errors.Is(err, domain.OrderCreationError) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, err
	}
	return &apiSchema.OrderID{OrderID: int64(orderID)}, nil
}
