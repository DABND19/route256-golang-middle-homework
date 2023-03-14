package main

import (
	"context"
	"log"
	"net"
	transationManager "route256/libs/transactor/postgresql"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	serviceAPI "route256/loms/internal/handlers/v1"
	"route256/loms/internal/middlewares"
	ordersRepository "route256/loms/internal/repository/postgresql/orders"
	stocksRepository "route256/loms/internal/repository/postgresql/stocks"
	apiSchema "route256/loms/pkg/lomsv1"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := pgxpool.Connect(ctx, config.Data.Postgres.DSN)
	if err != nil {
		log.Fatalln(err)
	}

	db := transationManager.New(pgPool)
	stocksRepo := stocksRepository.New(db)
	ordersRepo := ordersRepository.New(db)
	service := domain.New(db, ordersRepo, stocksRepo)

	lomsV1 := serviceAPI.New(service)
	apiSchema.RegisterLomsV1Server(s, lomsV1)

	log.Println("Server listen on", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Couldn't start a server:", err)
	}
}
