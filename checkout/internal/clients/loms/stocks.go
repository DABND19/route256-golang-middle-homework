package loms

import (
	"context"
	"route256/checkout/internal/models"
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"github.com/pkg/errors"
)

func (c *Client) Stocks(ctx context.Context, sku models.SKU) ([]models.Stock, error) {
	reqPayload := &lomsServiceAPI.SKU{Sku: uint32(sku)}
	resPayload, err := c.lomsServiceClient.Stocks(ctx, reqPayload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to request stocks")
	}

	stocks := make([]models.Stock, 0, len(resPayload.Stocks))
	for _, stock := range resPayload.GetStocks() {
		stocks = append(stocks, models.Stock{
			WarehouseID: models.WarehouseID(stock.WarehouseID),
			Count:       stock.Count,
		})
	}
	return stocks, nil
}
