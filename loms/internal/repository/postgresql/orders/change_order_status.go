package order

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
)

func toPostgresOrderStatus(status models.OrderStatus) (string, error) {
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

func (repo *Repository) ChangeOrderStatus(
	ctx context.Context,
	orderID models.OrderID,
	status models.OrderStatus,
) error {
	db := repo.GetEngine(ctx)

	postgresOrderStatus, err := toPostgresOrderStatus(status)
	if err != nil {
		return err
	}

	q := repo.queryBuilder.Update(ordersTableName).
		Set("status", postgresOrderStatus).
		Where(sq.Eq{"order_id": orderID})

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	res, err := db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domain.OrderNotFoundError
	}

	return nil
}
