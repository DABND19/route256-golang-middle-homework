package main

import (
	"context"
	"os"
	"os/signal"
	"route256/libs/logger"
	"route256/loms/internal/app"

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
