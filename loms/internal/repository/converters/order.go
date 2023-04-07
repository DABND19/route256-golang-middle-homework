package converters

import (
	"fmt"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schemas"
)

func ToDatabaseOrderStatus(status models.OrderStatus) (string, error) {
	orderStatusMapping := map[models.OrderStatus]string{
		models.OrderStatusNew:             schemas.OrderStatusNew,
		models.OrderStatusAwaitingPayment: schemas.OrderStatusAwaitingPayment,
		models.OrderStatusFailed:          schemas.OrderStatusFailed,
		models.OrderStatusPayed:           schemas.OrderStatusPayed,
		models.OrderStatusCancelled:       schemas.OrderStatusCancelled,
	}
	convertedValue, ok := orderStatusMapping[status]
	if !ok {
		return "", fmt.Errorf("Invalid order status value %s", string(status))
	}
	return convertedValue, nil
}

func ToDomainOrderStatus(status string) (models.OrderStatus, error) {
	orderStatusMapping := map[string]models.OrderStatus{
		schemas.OrderStatusNew:             models.OrderStatusNew,
		schemas.OrderStatusAwaitingPayment: models.OrderStatusAwaitingPayment,
		schemas.OrderStatusFailed:          models.OrderStatusFailed,
		schemas.OrderStatusPayed:           models.OrderStatusPayed,
		schemas.OrderStatusCancelled:       models.OrderStatusCancelled,
	}
	convertedValue, ok := orderStatusMapping[status]
	if !ok {
		return "", fmt.Errorf("Invalid order_status enum value %s", status)
	}
	return convertedValue, nil
}
