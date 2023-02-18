package createorder

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type RequestPayload struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}

type ResponsePayload struct {
	OrderID int64 `json:"orderID"`
}

func (*Handler) Handle(req RequestPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{OrderID: 1}
	return resPayload, nil
}

func (RequestPayload) Validate() error {
	return nil
}
