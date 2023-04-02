package notifications

import (
	"route256/loms/internal/domain"

	"github.com/Shopify/sarama"
)

type WorkerPool interface {
	Submit(task func())
}

type Client struct {
	producer             sarama.SyncProducer
	orderStatusTopicName string
	workerPool           WorkerPool
}

func New(
	syncProducer sarama.SyncProducer,
	orderStatusTopicName string,
	workerPool WorkerPool,
) domain.NotificationsClient {
	return &Client{
		producer:             syncProducer,
		orderStatusTopicName: orderStatusTopicName,
		workerPool:           workerPool,
	}
}
