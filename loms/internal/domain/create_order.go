package domain

import (
	"context"
	"errors"
	"log"
	"route256/loms/internal/models"
	"sort"
	"time"
)

func (s *Service) reserveItem(ctx context.Context, orderID models.OrderID, item models.OrderItem) error {
	stocks, err := s.StocksRespository.GetStocks(ctx, item.SKU)
	if err != nil {
		return err
	}

	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Count > stocks[j].Count
	})

	total := uint64(item.Count)
	for _, stock := range stocks {
		var reserveCount uint64
		if total <= stock.Count {
			reserveCount = total
		} else {
			reserveCount = stock.Count
		}

		err = s.StocksRespository.CreateItemBooking(ctx, orderID, stock.WarehouseID, item.SKU, uint16(reserveCount))
		if err != nil {
			return err
		}

		err = s.StocksRespository.UpdateStockItemsCount(ctx, stock.WarehouseID, item.SKU, -int64(reserveCount))
		if err != nil {
			return err
		}

		total -= reserveCount
		if total == 0 {
			break
		}
	}

	if total > 0 {
		return InsufficientStocksError
	}

	return nil
}

func (s *Service) CreateOrder(
	ctx context.Context,
	user models.User,
	items []models.OrderItem,
) (*models.OrderID, error) {
	var orderID *models.OrderID
	var failedReservationError error
	err := s.RunSerializable(ctx, func(ctx context.Context) error {
		var err error
		orderID, err = s.OrdersRespository.CreateOrder(ctx, user, items)
		if err != nil {
			return err
		}

		failedReservationError = s.RunInSavepoint(ctx, func(ctx context.Context) error {
			for _, item := range items {
				err = s.reserveItem(ctx, *orderID, item)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if failedReservationError != nil && !errors.Is(failedReservationError, InsufficientStocksError) {
			return failedReservationError
		}

		newOrderStatus := models.OrderStatusAwaitingPayment
		if failedReservationError != nil {
			newOrderStatus = models.OrderStatusFailed
		}

		err = s.OrdersRespository.ChangeOrderStatus(ctx, *orderID, newOrderStatus)
		if err != nil {
			return err
		}

		return nil
	})
	if failedReservationError != nil {
		return nil, failedReservationError
	}

	s.cancelOrderScheduler.Schedule(time.Now().Add(s.unpaidOrderTtl), func() {
		// Для того, чтобы никто не смог отменить задачку извне, создаем новый контекст
		err := s.CancelOrder(context.Background(), *orderID)
		if errors.Is(err, OrderAlreadyPayedError) {
			return
		}
		if err != nil {
			log.Printf("Failed to cancel order #%d, reason %s\n", *orderID, err)
			return
		}
		log.Printf("Order #%d cancelled\n", *orderID)
	})

	return orderID, err
}
