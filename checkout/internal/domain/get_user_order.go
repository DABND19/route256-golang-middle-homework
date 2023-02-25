package domain

import "context"

func (s *Service) GetUserOrder(ctx context.Context, user int64) ([]OrderItem, error) {
	return []OrderItem{
		{SKU: 1076963, Count: 1},
		{SKU: 1148162, Count: 1},
		{SKU: 1625903, Count: 1},
		{SKU: 2618151, Count: 1},
		{SKU: 2956315, Count: 1},
		{SKU: 2958025, Count: 1},
		{SKU: 3596599, Count: 1},
		{SKU: 3618852, Count: 1},
		{SKU: 4288068, Count: 1},
		{SKU: 4465995, Count: 1},
	}, nil
}
