package deletefromcart

import "context"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type ResponsePayload struct{}

func (*Handler) Handle(ctx context.Context, reqPayload RequestPayload) (ResponsePayload, error) {
	return ResponsePayload{}, nil
}

func (RequestPayload) Validate() error {
	return nil
}
