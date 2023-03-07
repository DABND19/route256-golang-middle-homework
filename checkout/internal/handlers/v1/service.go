package v1

import (
	"route256/checkout/internal/domain"
	apiSchema "route256/checkout/pkg/checkoutv1"
)

type Service struct {
	apiSchema.UnimplementedCheckoutV1Server
	service *domain.Service
}

func New(service *domain.Service) *Service {
	return &Service{
		apiSchema.UnimplementedCheckoutV1Server{},
		service,
	}
}
