package carts

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) UpdateCartItemProductsCount(
	ctx context.Context,
	user models.User,
	sku models.SKU,
	diff int32,
) (int32, error) {
	db := r.GetEngine(ctx)

	q := r.queryBuilder.Update(cartsItemsTableName).
		Set("count", sq.Expr("count + ?", diff)).
		Where(sq.Eq{"user_id": user}).
		Where(sq.Eq{"sku": sku}).
		Suffix("RETURNING count")

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	var newCount int32
	err = db.QueryRow(ctx, sql, args...).Scan(&newCount)
	if err != nil {
		if pgxscan.NotFound(err) {
			return 0, domain.CartItemNotFoundError
		}
		return 0, err
	}

	return newCount, nil
}
