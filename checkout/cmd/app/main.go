package main

import (
	"context"
	"os"
	"os/signal"
	"route256/checkout/internal/app"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	application := app.New(ctx)
	if err := application.Run(ctx); err != nil {
		logger.Fatal("Couldn't run the application.", zap.Error(err))
	}
	logger.Info("Application stopped.")
}
