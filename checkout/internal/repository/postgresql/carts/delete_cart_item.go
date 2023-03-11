package carts

import (
	"context"
	"route256/checkout/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteCartItem(
	ctx context.Context,
	user models.User,
	sku models.SKU,
) error {
	db := r.GetEngine(ctx)

	q := r.queryBuilder.Delete(cartsItemsTableName).
		Where(sq.Eq{"user_id": user}).
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
