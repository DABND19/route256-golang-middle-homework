package stocks

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *Repository) UpdateStockItemsCount(
	ctx context.Context,
	warehouseID models.WarehouseID,
	sku models.SKU,
	diff int64,
) error {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Update(itemsStocksTableName).
		Set("count", sq.Expr("count + ?", diff)).
		Where(sq.Eq{"warehouse_id": warehouseID}).
		Where(sq.Eq{"sku": sku}).
		Suffix("RETURNING count")

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	var balance int64
	err = db.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		if pgxscan.NotFound(err) {
			return domain.StockNotFoundError
		}
		return err
	}
	if balance < 0 {
		return domain.InsufficientStocksError
	}

	return nil
}
