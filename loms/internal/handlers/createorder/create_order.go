package createorder

import (
	"context"
	"route256/loms/internal/schemas"

	"github.com/pkg/errors"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	User  int64                      `json:"user"`
	Items []schemas.OrderItemPayload `json:"items"`
}

func (*Handler) Handle(ctx context.Context, req RequestPayload) (schemas.OrderPayload, error) {
	resPayload := schemas.OrderPayload{OrderID: 1}
	return resPayload, nil
}

func (p RequestPayload) Validate() error {
	if p.User == 0 {
		return errors.New("user required")
	}
	if len(p.Items) == 0 {
		return errors.New("empty order")
	}
	for pos, item := range p.Items {
		if err := item.Validate(); err != nil {
			return errors.Wrapf(err, "Invalid #%d item", pos)
		}
	}
	return nil
}
