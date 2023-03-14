package order

import (
	transactor "route256/libs/transactor/postgresql"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	transactor.EngineProvider
	queryBuilder sq.StatementBuilderType
}

const (
	ordersTableName      = "orders"
	ordersItemsTableName = "orders_items"
)

var (
	ordersColumns      = []string{"order_id", "user_id", "status"}
	ordersItemsColumns = []string{"order_id", "sku", "count"}
)

func New(engineProvider transactor.EngineProvider) domain.OrdersRespository {
	return &Repository{engineProvider, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
