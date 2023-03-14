package order

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schemas"
)

func (repo *Repository) queryCreateOrder(ctx context.Context, user models.User) (*models.OrderID, error) {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Insert(ordersTableName).
		Columns("user_id", "status").
		Values(user, schemas.OrderStatusNew).
		Suffix("RETURNING order_id")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var orderID models.OrderID
	err = db.QueryRow(ctx, sql, args...).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	return &orderID, nil
}

func (repo *Repository) queryCreateOrderItems(ctx context.Context, orderID models.OrderID, items []models.OrderItem) error {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Insert(ordersItemsTableName).Columns(ordersItemsColumns...)
	for _, orderItem := range items {
		q = q.Values(orderID, orderItem.SKU, orderItem.Count)
	}
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

func (repo *Repository) CreateOrder(
	ctx context.Context,
	user models.User,
	items []models.OrderItem,
) (*models.OrderID, error) {
	orderID, err := repo.queryCreateOrder(ctx, user)
	if err != nil {
		return nil, err
	}

	err = repo.queryCreateOrderItems(ctx, *orderID, items)
	if err != nil {
		return nil, err
	}

	return orderID, nil
}
