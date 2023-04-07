package main

import (
	"context"
	"os"
	"os/signal"
	"route256/libs/consumerwrapper"
	"route256/libs/logger"
	offsetstorage "route256/libs/offsetstorage/mock"
	"route256/notifications/internal/config"
	"route256/notifications/internal/handlers"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func NewConsumer(brokers []string) (sarama.Consumer, error) {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = false
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		return nil, err
	}
	return consumer, err
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger.Init(false)

	err := config.Load("config.yml")
	if err != nil {
		logger.Fatal("Failed to load config.", zap.Error(err))
	}

	consumer, err := NewConsumer(config.Data.Brokers)
	if err != nil {
		logger.Fatal("Couldn't connect to kafka cluster.", zap.Error(err))
	}
	offsetStorage := offsetstorage.New()
	wrapper := consumerwrapper.New(consumer, offsetStorage)

	errorsChan, err := wrapper.Subscribe(
		ctx,
		config.Data.OrderStatusChangeNotificationsTopicName,
		handlers.HandleOrderStatusChange,
	)
	if err != nil {
		logger.Fatal("Failed to subscribe to topic.", zap.Error(err))
	}
	go func() {
		for errorMessage := range errorsChan {
			logger.Error(
				"Order status change notification error.",
				zap.Error(errorMessage.Err),
				zap.Binary("messageKey", errorMessage.Message.Key),
				zap.Binary("messageValue", errorMessage.Message.Value),
			)
		}
	}()

	<-ctx.Done()
}
