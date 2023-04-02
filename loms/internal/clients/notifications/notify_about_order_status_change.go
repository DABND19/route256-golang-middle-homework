package notifications

import (
	"context"
	"log"
	"route256/loms/internal/models"
	"route256/notifications/pkg/notificationsv1"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

func (c *Client) NotifyAboutOrderStatusChange(ctx context.Context, orderID models.OrderID, orderStatus models.OrderStatus) {
	payload := notificationsv1.OrderStatusChangeNotification{
		OrderID: int64(orderID),
		Status:  string(orderStatus),
	}
	rawPayload, err := proto.Marshal(&payload)
	if err != nil {
		log.Println("Failed to serialize payload:", err)
		return
	}

	mesage := sarama.ProducerMessage{
		Topic:     c.orderStatusTopicName,
		Key:       sarama.StringEncoder(strconv.FormatInt(int64(orderID), 10)),
		Value:     sarama.ByteEncoder(rawPayload),
		Timestamp: time.Now(),
	}

	sendErrors := make(chan error)
	task := func() {
		_, _, err = c.producer.SendMessage(&mesage)
		if err == nil {
			close(sendErrors)
			return
		}

		select {
		case sendErrors <- err:
		case <-ctx.Done():
		}
	}
	c.workerPool.Submit(task)
	go func() {
		for {
			select {
			case err, ok := <-sendErrors:
				if !ok {
					return
				}
				log.Println("Failed to send message:", err)
				c.workerPool.Submit(task)
			case <-ctx.Done():
				return
			}
		}
	}()
}
