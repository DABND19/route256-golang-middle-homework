package cachemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HitsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "route256",
			Subsystem: "cache",
			Name:      "hits_count",
		},
		[]string{"target_name"},
	)
	MissesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "route256",
			Subsystem: "cache",
			Name:      "misses_count",
		},
		[]string{"target_name"},
	)
	HistogramReadTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "route256",
			Subsystem: "cache",
			Name:      "histogram_read_time_seconds",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		},
		[]string{"target_name"},
	)
	HistogramWriteTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "route256",
			Subsystem: "cache",
			Name:      "histogram_write_time_seconds",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		},
		[]string{"target_name"},
	)
)
