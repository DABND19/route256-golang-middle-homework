package ordersstatuschanges

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/converters"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) MarkChangeAsSubmitted(
	ctx context.Context,
	orderStatusChange models.OrderStatusChange,
) error {
	db := r.GetEngine(ctx)

	convertedStatus, err := converters.ToDatabaseOrderStatus(orderStatusChange.Status)
	if err != nil {
		return nil
	}

	q := r.queryBuilder.Update(ordersStatusChangesTableName).
		Set("submitted_at", time.Now()).
		Where(sq.Eq{"order_id": orderStatusChange.OrderID}).
		Where(sq.Eq{"status": convertedStatus})

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
