package notifications

import (
	"context"
	"route256/loms/internal/models"
	apiSchema "route256/loms/pkg/lomsv1"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

func (c *Client) NotifyAboutOrderStatusChange(ctx context.Context, change models.OrderStatusChange) error {
	payload := apiSchema.OrderStatusChange{
		OrderID:       int64(change.OrderID),
		UpdatedStatus: string(change.Status),
	}
	rawPayload, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}

	mesage := sarama.ProducerMessage{
		Topic:     c.orderStatusTopicName,
		Key:       sarama.StringEncoder(strconv.FormatInt(int64(change.OrderID), 10)),
		Value:     sarama.ByteEncoder(rawPayload),
		Timestamp: time.Now(),
	}

	_, _, err = c.producer.SendMessage(&mesage)
	if err != nil {
		return err
	}

	return nil
}
