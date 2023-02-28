package v1

import (
	"context"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) Stocks(ctx context.Context, reqPayload *apiSchema.SKU) (*apiSchema.StocksList, error) {
	return &apiSchema.StocksList{
		Stocks: []*apiSchema.Stock{
			{
				WarehouseID: 1,
				Count:       1,
			},
			{
				WarehouseID: 2,
				Count:       2,
			},
		},
	}, nil
}
