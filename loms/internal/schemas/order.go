package schemas

import "errors"

var (
	ErrMissedOrderID = errors.New("orderID required")
)

type OrderPayload struct {
	OrderID int64 `json:"orderID"`
}

func (p OrderPayload) Validate() error {
	if p.OrderID == 0 {
		return ErrMissedOrderID
	}
	return nil
}
