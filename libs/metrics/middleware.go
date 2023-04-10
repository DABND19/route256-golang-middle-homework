package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ServerMetricsMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	RequestsCounter.Inc()

	timeStart := time.Now()
	res, err := handler(ctx, req)

	resStatus := status.Code(err)
	ResponseCounter.WithLabelValues(info.FullMethod, resStatus.String()).Inc()

	reqDuration := time.Since(timeStart)
	HistogramResponseTime.WithLabelValues(info.FullMethod, resStatus.String()).Observe(reqDuration.Seconds())

	return res, err
}
