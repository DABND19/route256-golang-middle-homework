package main

import (
	"log"
	"net/http"

	"route256/libs/serverwrapper"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	createOrderHandler := createorder.New()
	listOrderHandler := listorder.New()
	orderPayedHandler := orderpayed.New()
	cancelOrderHandler := cancelorder.New()
	stocksHandler := stocks.New()

	http.Handle("/createOrder", serverwrapper.New(createOrderHandler.Handle))
	http.Handle("/listOrder", serverwrapper.New(listOrderHandler.Handle))
	http.Handle("/orderPayed", serverwrapper.New(orderPayedHandler.Handle))
	http.Handle("/cancelOrder", serverwrapper.New(cancelOrderHandler.Handle))
	http.Handle("/stocks", serverwrapper.New(stocksHandler.Handle))

	log.Println("Starting a server...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Couldn't start a server:", err)
		return
	}
}
