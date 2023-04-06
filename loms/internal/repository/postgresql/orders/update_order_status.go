package order

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/converters"

	sq "github.com/Masterminds/squirrel"
)

func (repo *Repository) UpdateOrderStatus(
	ctx context.Context,
	orderID models.OrderID,
	status models.OrderStatus,
) error {
	db := repo.GetEngine(ctx)

	postgresOrderStatus, err := converters.ToDatabaseOrderStatus(status)
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
