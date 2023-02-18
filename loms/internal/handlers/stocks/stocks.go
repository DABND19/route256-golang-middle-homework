package stocks

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type RequestPayload struct {
	SKU uint32 `json:"sku"`
}

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type ResponsePayload struct {
	Stocks []Stock `json:"stocks"`
}

func (*Handler) Handle(reqPayload RequestPayload) (ResponsePayload, error) {
	resPayload := ResponsePayload{
		Stocks: []Stock{
			{WarehouseID: 1, Count: 1},
			{WarehouseID: 2, Count: 2},
		},
	}
	return resPayload, nil
}

func (RequestPayload) Validate() error {
	return nil
}
