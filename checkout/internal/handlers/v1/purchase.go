package v1

import (
	"context"
	"route256/checkout/internal/models"
	apiSchema "route256/checkout/pkg/checkoutv1"
)

func (s *Service) Purchase(ctx context.Context, reqPayload *apiSchema.User) (*apiSchema.OrderID, error) {
	orderID, err := s.service.MakePurchase(ctx, models.User(reqPayload.User))
	if err != nil {
		return nil, err
	}
	return &apiSchema.OrderID{OrderID: int64(*orderID)}, nil
}
