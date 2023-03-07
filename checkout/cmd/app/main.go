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
	)
	if err != nil {
		log.Fatalln("Couldn't connect to product service:", err)
	}

	service := domain.New(lomsServiceClient, productServiceClient)

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
