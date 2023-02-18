package main

import (
	"log"
	"net/http"

	"route256/libs/serverwrapper"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/orderpayed"
)

const port = ":8081"

func main() {
	createOrderHandler := createorder.New()
	orderPayedHandler := orderpayed.New()

	http.Handle("/createOrder", serverwrapper.New(createOrderHandler.Handle))
	http.Handle("/orderPayed", serverwrapper.New(orderPayedHandler.Handle))

	log.Println("Starting a server...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Couldn't start a server:", err)
		return
	}
}
