package listcart

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/schemas"
)

type Handler struct {
	service *domain.Service
}

func New(service *domain.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type ItemPayload struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type ResponsePayload struct {
	Items      []ItemPayload `json:"items"`
	TotalPrice uint32        `json:"totalPrice"`
}

func (h *Handler) Handle(ctx context.Context, reqPayload schemas.UserPayload) (ResponsePayload, error) {
	cart, err := h.service.ListCart(ctx, reqPayload.User)
	if err != nil {
		return ResponsePayload{}, err
	}
	totalPrice := h.service.CalculateTotalPrice(cart)

	resPayload := ResponsePayload{
		TotalPrice: totalPrice,
		Items:      make([]ItemPayload, 0, len(cart)),
	}
	for _, item := range cart {
		resPayload.Items = append(resPayload.Items, ItemPayload{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	return resPayload, nil
}
