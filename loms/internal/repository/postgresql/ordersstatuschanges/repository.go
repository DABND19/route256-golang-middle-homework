package ordersstatuschanges

import (
	transactor "route256/libs/transactor/postgresql"
	"route256/loms/internal/domain"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	transactor.EngineProvider
	queryBuilder sq.StatementBuilderType
}

const ordersStatusChangesTableName = "orders_status_changes"

var (
	ordersStatusChangesColumns = []string{"created_at", "submitted_at", "order_id", "status"}
)

func New(engineProvider transactor.EngineProvider) domain.OrderStatusChangeRepository {
	return &Repository{engineProvider, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
