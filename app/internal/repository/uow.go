package repository

import (
	"JoinUp/internal/exceptions"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

type UOW struct {
	pool *pgxpool.Pool
}

type uowKey struct{}

type sqlRun interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewUOW(pool *pgxpool.Pool) UOW {
	return UOW{pool: pool}
}

func (t *UOW) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, uowKey{}, tx), nil
}

func (t *UOW) Commit(ctx context.Context) error {
	tx, err := t.txFromCtx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (t *UOW) Rollback(ctx context.Context) error {
	tx, err := t.txFromCtx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback(ctx)
}

// возвращает ошибку, так как если мы вызываем Rollback или Commit, значит явно хотим использовать транзакцию
func (t *UOW) txFromCtx(ctx context.Context) (pgx.Tx, error) {
	v := ctx.Value(uowKey{})
	if v == nil {
		return nil, exceptions.ErrNoTxObj
	}

	tx, ok := v.(pgx.Tx)
	if !ok {
		return nil, exceptions.ErrTxType
	}

	return tx, nil
}

// разница с одноименным методом у UOW в том, что если объект по ключу не найден, то это не ошибка.
// эта функция предназначечна для использования внутри репозиториев. Если в контексте нет объекта транзакции, значит запрос нужно выполнять без привязки к внешней транзакции
func TxFromCtx(ctx context.Context) (pgx.Tx, error) {
	v := ctx.Value(uowKey{})
	if v == nil {
		return nil, nil
	}

	tx, ok := v.(pgx.Tx)
	if !ok {
		return nil, exceptions.ErrTxType
	}

	return tx, nil
}
