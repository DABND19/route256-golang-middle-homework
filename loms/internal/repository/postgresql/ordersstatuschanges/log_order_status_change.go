package ordersstatuschanges

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/converters"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) LogOrderStatusChange(
	ctx context.Context,
	orderID models.OrderID,
	status models.OrderStatus,
) error {
	db := r.GetEngine(ctx)

	convertedStatus, err := converters.ToDatabaseOrderStatus(status)
	if err != nil {
		return err
	}

	q := r.queryBuilder.Insert(ordersStatusChangesTableName).
		Columns(ordersStatusChangesColumns...).
		Values(sq.Expr("NOW()"), nil, orderID, convertedStatus)

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
