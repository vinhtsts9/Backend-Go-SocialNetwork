package database

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	CreateCommentTx(ctx context.Context, arg CreateCommentTxParams) (CreateCommentTxResult, error)
}
type SQLStore struct {
	connPool *sql.DB
	*Queries
}

func NewStore(connPool *sql.DB) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
