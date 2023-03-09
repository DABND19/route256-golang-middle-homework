package stocks

import (
	"context"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schemas"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *Repository) GetItemBookings(
	ctx context.Context,
	orderID models.OrderID,
	sku models.SKU,
) ([]models.ItemBooking, error) {
	db := repo.GetEngine(ctx)

	q := repo.queryBuilder.Select(itemsBookingsColumns...).
		From(itemsBookingsTableName).
		Where(sq.Eq{"order_id": orderID}).
		Where(sq.Eq{"sku": sku})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	bookingRecords := make([]schemas.ItemBooking, 0)
	err = pgxscan.Select(ctx, db, &bookingRecords, sql, args...)
	if err != nil {
		return nil, err
	}

	return toItemBookingListModel(bookingRecords), nil
}

func toItemBookingListModel(bookingRecords []schemas.ItemBooking) []models.ItemBooking {
	bookings := make([]models.ItemBooking, 0, len(bookingRecords))
	for _, record := range bookingRecords {
		bookings = append(bookings, models.ItemBooking{
			WarehouseID: models.WarehouseID(record.WarehouseID),
			Count:       record.Count,
		})
	}
	return bookings
}
