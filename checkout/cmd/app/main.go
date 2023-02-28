package main

import (
	"errors"
	"log"
	"net/http"
	"route256/checkout/internal/clients/getproduct"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/serverwrapper"
	"route256/libs/serviceclient"

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

	lomsServiceClient := loms.New(lomsServiceConn)

	productServiceClient := serviceclient.New(config.Data.ExternalServices.Product.Url)
	getProductEndpointClient := getproduct.New(
		productServiceClient,
		"/get_product",
		config.Data.ExternalServices.Product.AccessToken,
	)

	service := domain.New(
		lomsServiceClient,
		getProductEndpointClient,
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
