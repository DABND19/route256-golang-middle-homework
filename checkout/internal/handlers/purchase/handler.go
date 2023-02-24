package purchase

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/schemas"
)

type Handler struct {
	service *domain.Service
}

func New(service *domain.Service) *Handler {
	return &Handler{service: service}
}

type ResponsePayload struct {
	OrderID int64 `json:"orderID"`
}

func (h *Handler) Handle(ctx context.Context, reqPayload schemas.UserPayload) (ResponsePayload, error) {
	orderID, err := h.service.MakePurchase(ctx, reqPayload.User)
	return ResponsePayload{OrderID: int64(orderID)}, err
}
