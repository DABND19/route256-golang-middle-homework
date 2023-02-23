package addtocart

import (
	"context"
	"errors"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/libs/serverwrapper"
)

type Handler struct {
	service *domain.Service
}

func New(service *domain.Service) *Handler {
	return &Handler{service: service}
}

type RequestPayload struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
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

func (RequestPayload) Validate() error {
	return nil
}
