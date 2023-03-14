package stocks

import (
	"context"
	"route256/loms/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (repo *Repository) DeleteItemBooking(
	ctx context.Context,
	orderID models.OrderID,
	warehouseID models.WarehouseID,
	sku models.SKU,
) error {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Delete(itemsBookingsTableName).
		Where(sq.Eq{"order_id": orderID}).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku})

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
