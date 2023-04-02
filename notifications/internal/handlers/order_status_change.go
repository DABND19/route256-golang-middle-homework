package handlers

import (
	"log"
	apiSchema "route256/notifications/pkg/notificationsv1"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

func HandleOrderStatusChange(message *sarama.ConsumerMessage) error {
	var payload apiSchema.OrderStatusChangeNotification
	err := proto.Unmarshal(message.Value, &payload)
	if err != nil {
		return err
	}

	log.Printf("Order #%d has changed status to %s\n", payload.OrderID, payload.Status)
	return nil
}
