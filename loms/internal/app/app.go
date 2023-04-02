package app

import (
	"context"
	"log"
	"net"
	"route256/loms/internal/config"
	serviceAPI "route256/loms/internal/handlers/v1"
	"route256/loms/internal/middlewares"
	apiSchema "route256/loms/pkg/lomsv1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
	if err := config.Load("config.yml"); err != nil {
		log.Fatalln("Failed to load app config:", err)
	}
	log.Println(config.Data)

	app.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(middlewares.DomainErrorsMiddleware),
		),
	)
	reflection.Register(app.grpcServer)

	app.dependencies = NewDependenciesProvider()
	lomsAPI := serviceAPI.New(app.dependencies.GetLOMSService(ctx))
	apiSchema.RegisterLomsV1Server(app.grpcServer, lomsAPI)
}

func (app *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", config.Data.Server.Address)
	if err != nil {
		return err
	}

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
