package tracing

import (
	"route256/libs/logger"

	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func Init(
	serviceName string,
	tracesCollectorEndpoint string,
) {
	config := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			CollectorEndpoint: tracesCollectorEndpoint,
		},
	}

	_, err := config.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Failed to init tracer.", zap.Error(err))
	}
}
