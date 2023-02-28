package main

import (
	"errors"
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/serverwrapper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	addToCartHandler := addtocart.New(service)
	deleteFromCartHandler := deletefromcart.New()
	listCartHandler := listcart.New(service)
	purchaseHandler := purchase.New(service)

	http.Handle("/addToCart", serverwrapper.New(addToCartHandler.Handle))
	http.Handle("/deleteFromCart", serverwrapper.New(deleteFromCartHandler.Handle))
	http.Handle("/listCart", serverwrapper.New(listCartHandler.Handle))
	http.Handle("/purchase", serverwrapper.New(purchaseHandler.Handle))

	log.Println("Starting a server...")
	err = http.ListenAndServe(config.Data.Server.Address, nil)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server stopped")
		} else {
			log.Fatalln("Couldn't start a server:", err)
		}
	}
}
