package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"route256/checkout/internal/config"
	serviceAPI "route256/checkout/internal/handlers/v1"
	"route256/checkout/internal/middlewares"
	apiSchema "route256/checkout/pkg/checkoutv1"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer   *grpc.Server
	dependencies *DependenciesProvider
}

func New(ctx context.Context) *App {
	app := &App{}
	app.bootstrap(ctx)
	return app
}

func (app *App) bootstrap(ctx context.Context) {
	logger.Init(false)

	if err := config.Load("config.yml"); err != nil {
		logger.Fatal("Failed to load app config.", zap.Error(err))
	}

	tracing.Init(
		"checkout_app",
		config.Data.Server.TracesCollectorEndpoint,
	)

	app.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				metrics.ServerMetricsMiddleware,
				middlewares.DomainErrorsMiddleware,
			),
		),
	)
	reflection.Register(app.grpcServer)

	app.dependencies = NewDependenciesProvider()
	lomsAPI := serviceAPI.New(app.dependencies.GetCheckoutService(ctx))
	apiSchema.RegisterCheckoutV1Server(app.grpcServer, lomsAPI)
}

func runMetricsServer(address string) {
	http.Handle("/metrics", metrics.New())
	go func() {
		err := http.ListenAndServe(address, nil)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Couldn't start metrics server.", zap.Error(err))
		}
	}()
}

func (app *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", config.Data.Server.Address)
	if err != nil {
		return err
	}

	runMetricsServer(config.Data.Server.MetricsAddress)

	go func() {
		<-ctx.Done()
		app.grpcServer.GracefulStop()
		app.dependencies.Close()
	}()

	if err := app.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
