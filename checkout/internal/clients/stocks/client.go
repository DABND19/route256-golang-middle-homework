package stocks

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/serviceclient"

	"github.com/pkg/errors"
)

type Client struct {
	serviceClient *serviceclient.ServiceClient
	endpointPath  string
}

func New(serviceClient *serviceclient.ServiceClient, endpointPath string) *Client {
	return &Client{
		serviceClient: serviceClient,
		endpointPath:  endpointPath,
	}
}

type RequestPayload struct {
	SKU uint32 `json:"sku"`
}

type StockPayload struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type ResponsePayload struct {
	Stocks []StockPayload `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	reqPayload := RequestPayload{SKU: sku}
	resPayload := ResponsePayload{}
	err := serviceclient.MakeRequest(ctx, c.serviceClient, c.endpointPath, reqPayload, &resPayload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to request stocks")
	}

	stocks := make([]domain.Stock, 0, len(resPayload.Stocks))
	for _, stock := range resPayload.Stocks {
		stocks = append(stocks, domain.Stock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}
	return stocks, nil
}
