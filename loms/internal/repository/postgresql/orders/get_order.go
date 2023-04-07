package order

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/converters"
	"route256/loms/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *Repository) queryOrder(ctx context.Context, orderID models.OrderID) (*schemas.Order, error) {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Select(ordersColumns...).
		From(ordersTableName).
		Where(sq.Eq{"order_id": orderID})
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var orderRecord schemas.Order
	err = pgxscan.Get(ctx, db, &orderRecord, sql, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, domain.OrderNotFoundError
		}
		return nil, err
	}
	return &orderRecord, nil
}

func (repo *Repository) queryOrderItems(ctx context.Context, orderID models.OrderID) ([]schemas.OrderItem, error) {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Select(ordersItemsColumns...).
		From(ordersItemsTableName).
		Where(sq.Eq{"order_id": orderID})
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	orderItemRecords := make([]schemas.OrderItem, 0)
	err = pgxscan.Select(ctx, db, &orderItemRecords, sql, args...)
	if err != nil {
		return nil, err
	}
	return orderItemRecords, nil
}

func (repo *Repository) GetOrder(
	ctx context.Context,
	orderID models.OrderID,
) (*models.Order, error) {
	orderRecord, err := repo.queryOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	orderItemRecords, err := repo.queryOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order, err := toOrderModel(*orderRecord, orderItemRecords)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func toOrderModel(orderRecord schemas.Order, orderItemRecords []schemas.OrderItem) (*models.Order, error) {
	orderStatus, err := converters.ToDomainOrderStatus(orderRecord.Status)
	if err != nil {
		return nil, err
	}

	orderItems := make([]models.OrderItem, 0, len(orderItemRecords))
	for _, record := range orderItemRecords {
		orderItems = append(orderItems, models.OrderItem{
			SKU:   models.SKU(record.SKU),
			Count: record.Count,
		})
	}

	return &models.Order{
		User:   models.User(orderRecord.UserID),
		Status: orderStatus,
		Items:  orderItems,
	}, nil
}
