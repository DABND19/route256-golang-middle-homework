package v1

import (
	"context"
	apiSchema "route256/checkout/pkg/checkoutv1"
)

func (s *Service) Purchase(ctx context.Context, reqPayload *apiSchema.User) (*apiSchema.OrderID, error) {
	orderID, err := s.service.MakePurchase(ctx, reqPayload.User)
	return &apiSchema.OrderID{OrderID: int64(orderID)}, err
}
