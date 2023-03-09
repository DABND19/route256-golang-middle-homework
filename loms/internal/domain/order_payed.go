package domain

import (
	"context"
	"route256/loms/internal/models"
)

func (s *Service) OrderPayed(ctx context.Context, orderID models.OrderID) error {
	return s.RunSerializable(ctx, func(ctx context.Context) error {
		order, err := s.OrdersRespository.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		for _, item := range order.Items {
			itemBookings, err := s.StocksRespository.GetItemBookings(ctx, orderID, item.SKU)
			if err != nil {
				return err
			}

			for _, booking := range itemBookings {
				err = s.StocksRespository.DeleteItemBooking(ctx, orderID, booking.WarehouseID, item.SKU)
				if err != nil {
					return err
				}
			}
		}

		err = s.OrdersRespository.ChangeOrderStatus(ctx, orderID, models.OrderStatusPayed)
		if err != nil {
			return err
		}

		return nil
	})
}
