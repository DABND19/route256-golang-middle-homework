package schemas

import "errors"

var (
	ErrMissedSKU   = errors.New("sku required")
	ErrMissedCount = errors.New("count required")
)

type CartItemPayload struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

func (p CartItemPayload) Validate() error {
	if p.SKU == 0 {
		return ErrMissedSKU
	}
	if p.Count == 0 {
		return ErrMissedCount
	}
	return nil
}
