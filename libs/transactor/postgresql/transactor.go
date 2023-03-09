package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type engineKey string

const key = engineKey("engine")

type Engine interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type EngineProvider interface {
	GetEngine(ctx context.Context) Engine
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (t *TransactionManager) runInTransaction(
	ctx context.Context,
	txOptions pgx.TxOptions,
	txFn func(ctx context.Context) error,
) error {
	tx, err := t.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	if err := txFn(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *TransactionManager) RunReadCommited(ctx context.Context, txFn func(ctx context.Context) error) error {
	return t.runInTransaction(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, txFn)
}

func (t *TransactionManager) RunRepeatableRead(ctx context.Context, txFn func(ctx context.Context) error) error {
	return t.runInTransaction(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, txFn)
}

func (t *TransactionManager) RunSerializable(ctx context.Context, txFn func(ctx context.Context) error) error {
	return t.runInTransaction(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, txFn)
}

func (t *TransactionManager) RunInSavepoint(ctx context.Context, txFn func(ctx context.Context) error) error {
	tx, ok := ctx.Value(key).(pgx.Tx)
	if !ok || tx == nil {
		return errors.New("The transaction is not begun")
	}

	savepoint, err := tx.Begin(ctx)
	if err != nil {
		return err
	}

	if err := txFn(ctx); err != nil {
		return multierr.Combine(err, savepoint.Rollback(ctx))
	}

	if err := savepoint.Commit(ctx); err != nil {
		return multierr.Combine(err, savepoint.Rollback(ctx))
	}

	return nil
}

func (t *TransactionManager) GetEngine(ctx context.Context) Engine {
	engine, ok := ctx.Value(key).(Engine)
	if !ok || engine == nil {
		return t.pool
	}
	return engine
}
