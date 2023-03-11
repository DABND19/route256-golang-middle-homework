package carts

import (
	"route256/checkout/internal/domain"
	transactor "route256/libs/transactor/postgresql"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	transactor.EngineProvider
	queryBuilder sq.StatementBuilderType
}

func New(tr transactor.EngineProvider) domain.CartsRepository {
	return &Repository{
		tr, sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

const cartsItemsTableName = "carts_items"

var (
	cartsItemsColumnNames = []string{"user_id", "created_at", "sku", "count"}
)
