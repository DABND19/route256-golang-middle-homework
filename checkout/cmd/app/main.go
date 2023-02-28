package main

import (
	"log"
	"net"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	serviceAPI "route256/checkout/internal/handlers/v1"
	apiSchema "route256/checkout/pkg/checkoutv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Load("config.yml")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	lomsServiceConn, err := grpc.Dial(
		config.Data.ExternalServices.Loms.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Couldn't connect to LOMS service:", err)
	}
	lomsServiceClient := loms.New(lomsServiceConn)

	productServiceConn, err := grpc.Dial(
		config.Data.ExternalServices.Product.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Couldn't connect to product service:", err)
	}
	productServiceClient := product.New(
		productServiceConn,
		config.Data.ExternalServices.Product.AccessToken,
	)

	service := domain.New(
		lomsServiceClient,
		productServiceClient,
		lomsServiceClient,
	)

	lis, err := net.Listen("tcp", config.Data.Server.Address)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	checkoutV1 := serviceAPI.New(service)
	apiSchema.RegisterCheckoutV1Server(s, checkoutV1)

	log.Println("Server listen on:", lis.Addr())
	err = s.Serve(lis)
	if err != nil {
		log.Fatalln("Couldn't start a server:", err)
	}
}
