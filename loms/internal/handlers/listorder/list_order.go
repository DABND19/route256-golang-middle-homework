package listorder

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	OrderID int64 `json:"orderID"`
}

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type ResponsePayload struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

func (*Handler) Handle(req RequestPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{
		Status: "new",
		User:   1,
		Items: []Item{
			{SKU: 1, Count: 1},
		},
	}
	return resPayload, nil
}

func (RequestPayload) Validate() error {
	return nil
}
