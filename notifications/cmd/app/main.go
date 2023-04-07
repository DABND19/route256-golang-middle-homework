package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"route256/libs/consumerwrapper"
	offsetstorage "route256/libs/offsetstorage/mock"
	"route256/notifications/internal/config"
	"route256/notifications/internal/handlers"

	"github.com/Shopify/sarama"
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

	err := config.Load("config.yml")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	consumer, err := NewConsumer(config.Data.Brokers)
	if err != nil {
		log.Fatalln("Couldn't connect to kafka cluster:", err)
	}
	offsetStorage := offsetstorage.New()
	wrapper := consumerwrapper.New(consumer, offsetStorage)

	errorsChan, err := wrapper.Subscribe(
		ctx,
		config.Data.OrderStatusChangeNotificationsTopicName,
		handlers.HandleOrderStatusChange,
	)
	if err != nil {
		log.Fatalln("Failed to subscribe to topic:", err)
	}
	go func() {
		for errorMessage := range errorsChan {
			fmt.Println("Received error:", errorMessage.Err, "\nMessage:", *errorMessage.Message)
		}
	}()

	<-ctx.Done()
}
