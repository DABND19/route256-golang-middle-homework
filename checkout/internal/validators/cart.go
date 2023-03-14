package validators

import "route256/checkout/internal/domain"

type ProductsCountGetter interface {
	GetCount() uint32
}

func ValidateProductsCount(payload ProductsCountGetter) error {
	if payload.GetCount() <= 0 {
		return domain.InvalidProductsCount
	}
	return nil
}
