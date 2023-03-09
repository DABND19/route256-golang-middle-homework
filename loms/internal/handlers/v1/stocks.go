package v1

import (
	"context"
	"route256/loms/internal/models"
	apiSchema "route256/loms/pkg/lomsv1"
)

func (s *Service) Stocks(ctx context.Context, reqPayload *apiSchema.SKU) (*apiSchema.StocksList, error) {
	stocks, err := s.service.GetStocks(ctx, models.SKU(reqPayload.Sku))
	if err != nil {
		return nil, err
	}

	resPayload := apiSchema.StocksList{
		Stocks: make([]*apiSchema.Stock, 0, len(stocks)),
	}
	for _, stock := range stocks {
		resPayload.Stocks = append(resPayload.Stocks, &apiSchema.Stock{
			WarehouseID: int64(stock.WarehouseID),
			Count:       stock.Count,
		})
	}
	return &resPayload, nil
}
