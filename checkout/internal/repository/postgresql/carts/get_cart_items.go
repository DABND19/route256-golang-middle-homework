package carts

import (
	"context"
	"route256/checkout/internal/models"
	"route256/checkout/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) GetCartItems(
	ctx context.Context,
	user models.User,
) ([]models.CartItem, error) {
	db := r.GetEngine(ctx)

	q := r.queryBuilder.Select(cartsItemsColumnNames...).
		From(cartsItemsTableName).
		Where(sq.Eq{"user_id": user}).
		OrderBy("created_at DESC")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var cartItemRecords []schemas.CartItem
	err = pgxscan.Select(ctx, db, &cartItemRecords, sql, args...)
	if err != nil {
		return nil, err
	}
	return toCartItemsListModel(cartItemRecords), nil
}

func toCartItemsListModel(cartItemRecords []schemas.CartItem) []models.CartItem {
	cartItems := make([]models.CartItem, 0, len(cartItemRecords))
	for _, record := range cartItemRecords {
		cartItems = append(cartItems, models.CartItem{
			SKU:   models.SKU(record.SKU),
			Count: models.ProductsCount(record.Count),
		})
	}
	return cartItems
}
