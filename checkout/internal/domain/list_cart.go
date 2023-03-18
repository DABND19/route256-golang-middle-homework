package domain

import (
	"context"
	"route256/checkout/internal/models"
	"sync"

	"github.com/pkg/errors"
)

var (
	ProductServiceRateLimitError = errors.New("Too many requests to product service")
)

type fetchProductTaskResult struct {
	SKU models.SKU
	*models.Product
}

func (s *Service) fetchProducts(
	ctx context.Context,
	skus []models.SKU,
) (map[models.SKU]*models.Product, error) {
	fetchingCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	errorsQueue := make(chan error, len(skus))
	responsesQueue := make(chan fetchProductTaskResult, len(skus))
	for _, sku := range skus {
		sku := sku
		wg.Add(1)
		s.listCartWp.Submit(func() {
			product, err := s.productServiceClient.GetProduct(fetchingCtx, sku)
			if err != nil {
				errorsQueue <- err
				return
			}
			responsesQueue <- fetchProductTaskResult{sku, product}
		})
	}

	var fetchingErr error
	products := make(map[models.SKU]*models.Product, len(skus))
	go func() {
	loop:
		for {
			select {
			case res, ok := <-responsesQueue:
				if !ok {
					break loop
				}
				wg.Done()
				products[res.SKU] = res.Product

			case err, ok := <-errorsQueue:
				if !ok {
					break loop
				}
				wg.Done()
				if errors.Is(err, ProductNotFound) {
					continue
				}
				if fetchingErr == nil {
					fetchingErr = err
					cancel()
				}
			}
		}
	}()

	wg.Wait()
	close(responsesQueue)
	close(errorsQueue)

	return products, fetchingErr
}

func (s *Service) ListCart(ctx context.Context, user models.User) ([]models.CartProduct, error) {
	var cartItems []models.CartItem
	err := s.RunReadCommited(ctx, func(ctx context.Context) error {
		var err error
		cartItems, err = s.cartsRepository.GetCartItems(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query user cart")
	}

	skus := make([]models.SKU, 0, len(cartItems))
	for _, item := range cartItems {
		skus = append(skus, item.SKU)
	}
	fetchedProducts, err := s.fetchProducts(ctx, skus)
	if err != nil {
		return nil, err
	}

	cartProducts := make([]models.CartProduct, 0, len(cartItems))
	for _, item := range cartItems {
		product, ok := fetchedProducts[item.SKU]
		if !ok {
			continue
		}
		cartProducts = append(cartProducts, models.CartProduct{
			CartItem: item,
			Product:  *product,
		})
	}
	return cartProducts, nil
}

func (s *Service) CalculateTotalPrice(cart []models.CartProduct) (total uint32) {
	for _, item := range cart {
		total += item.Price * uint32(item.Count)
	}
	return
}
