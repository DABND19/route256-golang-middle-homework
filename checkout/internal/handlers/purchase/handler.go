package purchase

import (
	"context"
	"route256/checkout/internal/domain"
)

type Handler struct {
	service *domain.Service
}

func New(service *domain.Service) *Handler {
	return &Handler{service: service}
}

type RequestPayload struct {
	User int64 `json:"user"`
}

type ResponsePayload struct{}

func (h *Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	err := h.service.MakePurchase(ctx, reqPayload.User)
	return ResponsePayload{}, err
}

func (RequestPayload) Validate() error {
	return nil
}
