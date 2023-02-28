package loms

import (
	"context"
	"log"
	"route256/checkout/internal/domain"
	lomsServiceAPI "route256/loms/pkg/lomsv1"

	"github.com/pkg/errors"
)

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	reqPayload := &lomsServiceAPI.SKU{Sku: sku}
	resPayload, err := c.lomsServiceClient.Stocks(ctx, reqPayload)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "Failed to request stocks")
	}

	stocks := make([]domain.Stock, 0, len(resPayload.Stocks))
	for _, stock := range resPayload.GetStocks() {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.GetWarehouseID(),
			Count:       stock.GetCount(),
		})
	}
	return stocks, nil
}
