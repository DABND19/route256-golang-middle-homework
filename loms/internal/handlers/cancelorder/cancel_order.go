package cancelorder

import (
	"context"
	"route256/loms/internal/schemas"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type ResponsePayload struct{}

func (*Handler) Handle(ctx context.Context, reqPayload schemas.OrderPayload) (ResponsePayload, error) {
	return ResponsePayload{}, nil
}
