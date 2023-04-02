package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"route256/loms/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	application := app.New(ctx)
	if err := application.Run(ctx); err != nil {
		log.Fatalln("Couldn't run the application:", err)
	}
	log.Println("Application stopped.")
}
