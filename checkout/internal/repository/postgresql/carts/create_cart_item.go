package carts

import (
	"context"
	"route256/checkout/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateCartItem(
	ctx context.Context,
	user models.User,
	sku models.SKU,
	count models.ProductsCount,
) error {
	db := r.GetEngine(ctx)

	q := r.queryBuilder.Insert(cartsItemsTableName).
		Columns(cartsItemsColumnNames...).
		Values(user, sq.Expr("NOW()"), sku, count).
		Suffix(`
			ON CONFLICT (user_id, sku) DO UPDATE
			SET count = EXCLUDED.count
		`)

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
