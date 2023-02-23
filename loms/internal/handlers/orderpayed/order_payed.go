package orderpayed

import "context"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	OrderID int64
}

type ResponsePayload struct{}

func (*Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{}
	return resPayload, nil
}

func (RequestPayload) Validate() error {
	return nil
}
