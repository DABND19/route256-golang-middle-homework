package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/createorder"
	"route256/checkout/internal/clients/getproduct"
	"route256/checkout/internal/clients/stocks"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/serverwrapper"
	"route256/libs/serviceclient"
)

func main() {
	err := config.Load("config.yml")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	lomsServiceClient := serviceclient.New(config.Data.ExternalServices.Loms.Url)
	stocksEndpointClient := stocks.New(lomsServiceClient, "/stocks")
	createOrderEndpointClient := createorder.New(lomsServiceClient, "/createOrder")

	productServiceClient := serviceclient.New(config.Data.ExternalServices.Product.Url)
	getProductEndpointClient := getproduct.New(
		productServiceClient,
		"/get_product",
		config.Data.ExternalServices.Product.AccessToken,
	)

	service := domain.New(
		stocksEndpointClient,
		getProductEndpointClient,
		createOrderEndpointClient,
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
		log.Fatal("Couldn't start a server:", err)
	}
}
