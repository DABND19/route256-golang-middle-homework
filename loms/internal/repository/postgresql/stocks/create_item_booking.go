package stocks

import (
	"context"
	"route256/loms/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (repo *Repository) CreateItemBooking(
	ctx context.Context,
	orderID models.OrderID,
	warehouseID models.WarehouseID,
	sku models.SKU,
	count uint16,
) error {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Insert(itemsBookingsTableName).
		Columns(itemsBookingsColumns...).
		Values(sq.Expr("NOW()"), orderID, warehouseID, sku, count)

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
