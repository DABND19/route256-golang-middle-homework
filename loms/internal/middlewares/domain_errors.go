package middlewares

import (
	"context"
	"errors"
	"log"
	"route256/loms/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DomainErrorsMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	res, err := handler(ctx, req)
	if err == nil {
		return res, nil
	}

	if errors.Is(err, domain.InsufficientStocksError) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if errors.Is(err, domain.OrderNotFoundError) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.StockNotFoundError) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	log.Println(err)
	return nil, status.Error(codes.Internal, "Internal server error")
}
