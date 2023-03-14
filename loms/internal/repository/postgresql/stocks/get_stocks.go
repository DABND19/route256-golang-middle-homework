package stocks

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *Repository) GetStocks(ctx context.Context, sku models.SKU) ([]models.Stock, error) {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Select(itemsStocksColumns...).
		From(itemsStocksTableName).
		Where(sq.Eq{"sku": sku}).
		Where(sq.Gt{"count": 0})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	stockRecords := make([]schemas.ItemStock, 0)
	err = pgxscan.Select(ctx, db, &stockRecords, sql, args...)
	if err != nil {
		return nil, err
	}

	return toStocksListModel(stockRecords), nil
}

func toStocksListModel(stockRecords []schemas.ItemStock) []models.Stock {
	stocks := make([]models.Stock, 0, len(stockRecords))
	for _, record := range stockRecords {
		stocks = append(stocks, models.Stock{
			WarehouseID: models.WarehouseID(record.WarehouseID),
			Count:       record.Count,
		})
	}
	return stocks
}
