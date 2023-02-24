package stocks

import (
	"context"
	"errors"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
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

func (*Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{
		Stocks: []StockPayload{
			{WarehouseID: 1, Count: 1},
			{WarehouseID: 2, Count: 2},
		},
	}
	return resPayload, nil
}

func (p RequestPayload) Validate() error {
	if p.SKU == 0 {
		return errors.New("sku required")
	}
	return nil
}
