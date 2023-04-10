package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "http",
		Name:      "requests_total",
	})
	ResponseCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "route256",
			Subsystem: "http",
			Name:      "responses_total",
		},
		[]string{"method", "status"},
	)
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "http",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"method", "status"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}
