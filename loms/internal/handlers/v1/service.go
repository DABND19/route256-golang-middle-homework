package v1

import apiSchema "route256/loms/pkg/lomsv1"

type Service struct {
	apiSchema.UnimplementedLomsV1Server
}

func New() *Service {
	return &Service{
		apiSchema.UnimplementedLomsV1Server{},
	}
}
