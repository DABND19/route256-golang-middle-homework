package handlers

import (
	"route256/libs/logger"
	apiSchema "route256/loms/pkg/lomsv1"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func HandleOrderStatusChange(message *sarama.ConsumerMessage) error {
	var payload apiSchema.OrderStatusChange
	err := proto.Unmarshal(message.Value, &payload)
	if err != nil {
		return err
	}

	logger.Info(
		"Order has changed status.",
		zap.Int64("orderID", payload.OrderID),
		zap.String("updatedStatus", payload.UpdatedStatus),
	)
	return nil
}
