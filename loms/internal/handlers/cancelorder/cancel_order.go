package cancelorder

import "context"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	OrderID int64 `json:"orderID"`
}

type ResponsePayload struct{}

func (*Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	return ResponsePayload{}, nil
}

func (RequestPayload) Validate() error {
	return nil
}
