package main

import (
	"log"
	"net"
	"route256/loms/internal/config"
	serviceAPI "route256/loms/internal/handlers/v1"
	apiSchema "route256/loms/pkg/lomsv1"

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

	s := grpc.NewServer()
	reflection.Register(s)

	lomsV1 := serviceAPI.New()
	apiSchema.RegisterLomsV1Server(s, lomsV1)

	log.Println("Server listen on", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Couldn't start a server:", err)
	}
}
