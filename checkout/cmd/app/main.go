package main

import (
	"context"
	"log"
	"net"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	serviceAPI "route256/checkout/internal/handlers/v1"
	"route256/checkout/internal/middlewares"
	cartsRepository "route256/checkout/internal/repository/postgresql/carts"
	apiSchema "route256/checkout/pkg/checkoutv1"
	transactionManager "route256/libs/transactor/postgresql"
	"route256/libs/workerpool"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Load("config.yml")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	lomsServiceClient, err := loms.New(config.Data.ExternalServices.Loms.Url)
	if err != nil {
		log.Fatalln("Couldn't connect to LOMS service:", err)
	}

	productServiceClient, err := product.New(
		config.Data.ExternalServices.Product.Url,
		config.Data.ExternalServices.Product.AccessToken,
		int(config.Data.ExternalServices.Product.RateLimit),
	)
	if err != nil {
		log.Fatalln("Couldn't connect to product service:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := pgxpool.Connect(ctx, config.Data.Postgres.DSN)
	if err != nil {
		log.Fatalln(err)
	}

	db := transactionManager.New(pgPool)
	cartsRepo := cartsRepository.New(db)
	service := domain.New(
		db,
		cartsRepo,
		lomsServiceClient,
		productServiceClient,
		workerpool.New(5),
	)

	lis, err := net.Listen("tcp", config.Data.Server.Address)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(middlewares.DomainErrorsMiddleware),
		),
	)
	reflection.Register(s)

	checkoutV1 := serviceAPI.New(service)
	apiSchema.RegisterCheckoutV1Server(s, checkoutV1)

	log.Println("Server listen on:", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("Couldn't start a server:", err)
	}
}
