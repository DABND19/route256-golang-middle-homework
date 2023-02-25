package addtocart

import (
	"context"
	"errors"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/schemas"
	"route256/libs/serverwrapper"
)

type Handler struct {
	service *domain.Service
}

func New(service *domain.Service) *Handler {
	return &Handler{service: service}
}

type RequestPayload struct {
	schemas.UserPayload
	schemas.CartItemPayload
}

type ResponsePayload struct{}

func (h *Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	err := h.service.AddToCart(ctx, reqPayload.User, reqPayload.SKU, reqPayload.Count)
	if errors.Is(err, domain.InsufficientStocksError) {
		return ResponsePayload{}, serverwrapper.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	return ResponsePayload{}, err
}

func (p RequestPayload) Validate() error {
	if err := p.CartItemPayload.Validate(); err != nil {
		return err
	}
	if err := p.UserPayload.Validate(); err != nil {
		return err
	}
	return nil
}
