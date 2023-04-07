package middlewares

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	"route256/libs/logger"

	"go.uber.org/zap"
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

	if errors.Is(err, domain.OrderCreationError) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if errors.Is(err, domain.ProductNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.InvalidProductsCount) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if errors.Is(err, domain.CartItemNotFoundError) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, domain.ProductServiceRateLimitError) {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	logger.Error("Internal server error.", zap.Error(err))
	return nil, status.Error(codes.Internal, "Internal server error")
}
