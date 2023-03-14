package stocks

import (
	transactor "route256/libs/transactor/postgresql"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	transactor.EngineProvider
	queryBuilder sq.StatementBuilderType
}

func New(engineProvider transactor.EngineProvider) domain.StocksRespository {
	return &Repository{engineProvider, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

const (
	itemsStocksTableName   = "items_stocks"
	itemsBookingsTableName = "items_bookings"
)

var (
	itemsStocksColumns   = []string{"warehouse_id", "sku", "count"}
	itemsBookingsColumns = []string{"created_at", "order_id", "warehouse_id", "sku", "count"}
)
