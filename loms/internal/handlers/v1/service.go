package v1

import (
	"route256/loms/internal/domain"
	apiSchema "route256/loms/pkg/lomsv1"
)

type Service struct {
	apiSchema.UnimplementedLomsV1Server

	service *domain.Service
}

func New(service *domain.Service) *Service {
	return &Service{
		apiSchema.UnimplementedLomsV1Server{},
		service,
	}
}
