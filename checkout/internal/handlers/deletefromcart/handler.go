package deletefromcart

import (
	"context"
	"route256/checkout/internal/schemas"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	schemas.UserPayload
	schemas.CartItemPayload
}

type ResponsePayload struct{}

func (*Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	return ResponsePayload{}, nil
}

func (p RequestPayload) Validate() error {
	if err := p.UserPayload.Validate(); err != nil {
		return err
	}
	if err := p.CartItemPayload.Validate(); err != nil {
		return err
	}
	return nil
}
