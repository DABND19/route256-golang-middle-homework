package consumerwrapper

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
)

type ErrorMessage struct {
	Err     error
	Message *sarama.ConsumerMessage
}

type Wrapper struct {
	consumer      sarama.Consumer
	offsetStorage OffsetsStorage
}

type OffsetsStorage interface {
	GetOffset(ctx context.Context, partition int32) (int64, error)
	SetOffset(ctx context.Context, partition int32, value int64) error
}

func New(
	consumer sarama.Consumer,
	offsetStorage OffsetsStorage,
) *Wrapper {
	return &Wrapper{
		consumer:      consumer,
		offsetStorage: offsetStorage,
	}
}

func (w *Wrapper) Subscribe(
	ctx context.Context,
	topicName string,
	handler func(msg *sarama.ConsumerMessage) error,
) (<-chan ErrorMessage, error) {
	partitionsList, err := w.consumer.Partitions(topicName)
	if err != nil {
		return nil, err
	}

	wg := new(sync.WaitGroup)
	errorsChan := make(chan ErrorMessage, len(partitionsList))
	defer func() {
		go func() {
			wg.Wait()
			close(errorsChan)
		}()
	}()
	for _, partition := range partitionsList {
		partition := partition

		initialOffset, err := w.offsetStorage.GetOffset(ctx, partition)
		if err != nil {
			return nil, err
		}

		pc, err := w.consumer.ConsumePartition(topicName, partition, initialOffset)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case message, ok := <-pc.Messages():
					if !ok {
						return
					}
					err := handler(message)
					if err != nil {
						errorsChan <- ErrorMessage{
							Err:     err,
							Message: message,
						}
					}
					w.offsetStorage.SetOffset(ctx, partition, message.Offset)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return errorsChan, nil
}
