package handlers

import (
	"log"
	apiSchema "route256/loms/pkg/lomsv1"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

func HandleOrderStatusChange(message *sarama.ConsumerMessage) error {
	var payload apiSchema.OrderStatusChange
	err := proto.Unmarshal(message.Value, &payload)
	if err != nil {
		return err
	}

	log.Printf("Order #%d has changed status to %s\n", payload.OrderID, payload.UpdatedStatus)
	return nil
}
