package listorder

import (
	"context"
	"route256/loms/internal/schemas"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type ResponsePayload struct {
	Status string                     `json:"status"`
	User   int64                      `json:"user"`
	Items  []schemas.OrderItemPayload `json:"items"`
}

func (*Handler) Handle(ctx context.Context, req schemas.OrderPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{
		Status: "new",
		User:   1,
		Items: []schemas.OrderItemPayload{
			{SKU: 1, Count: 1},
		},
	}
	return resPayload, nil
}
