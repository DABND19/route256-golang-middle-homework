package ordersstatuschanges

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/converters"
	"route256/loms/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) GetUnsubmittedChanges(ctx context.Context) ([]models.OrderStatusChange, error) {
	db := r.GetEngine(ctx)

	q := r.queryBuilder.Select(ordersStatusChangesColumns...).
		From(ordersStatusChangesTableName).
		Where(sq.Eq{"submitted_at": nil})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	changesRecords := make([]schemas.OrderStatusChange, 0)
	err = pgxscan.Select(ctx, db, &changesRecords, sql, args...)
	if err != nil {
		return nil, err
	}

	return toOrderStatusChangeListModel(changesRecords)
}

func toOrderStatusChangeListModel(records []schemas.OrderStatusChange) ([]models.OrderStatusChange, error) {
	res := make([]models.OrderStatusChange, 0, len(records))
	for _, item := range records {
		convertedStatus, err := converters.ToDomainOrderStatus(item.Status)
		if err != nil {
			return nil, err
		}

		res = append(res, models.OrderStatusChange{
			CreatedAt: item.CreatedAt,
			OrderID:   models.OrderID(item.OrderID),
			Status:    convertedStatus,
		})
	}
	return res, nil
}
