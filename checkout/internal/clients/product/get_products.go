package product

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/models"
	"sync"
)

type taskResult[ID, Payload any] struct {
	ID      ID
	Payload Payload
	Error   error
}

func (c *Client) GetProducts(ctx context.Context, skus []models.SKU) (map[models.SKU]*models.Product, error) {
	fetchingCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	queue := make(chan taskResult[models.SKU, *models.Product])
	for _, sku := range skus {
		sku := sku
		wg.Add(1)
		c.wp.Submit(func() {
			product, err := c.GetProduct(fetchingCtx, sku)
			queue <- taskResult[models.SKU, *models.Product]{
				ID:      sku,
				Payload: product,
				Error:   err,
			}
		})
	}

	var fetchingErr error
	products := make(map[models.SKU]*models.Product, len(skus))
	go func() {
		for res := range queue {
			wg.Done()
			if errors.Is(res.Error, domain.ProductNotFound) {
				continue
			}
			if res.Error != nil && fetchingErr == nil {
				cancel()
				fetchingErr = res.Error
			}
			products[res.ID] = res.Payload
		}
	}()

	wg.Wait()
	close(queue)

	if fetchingErr != nil {
		return nil, fetchingErr
	}

	return products, nil
}
