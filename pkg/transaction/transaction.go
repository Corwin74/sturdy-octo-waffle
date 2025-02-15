package transaction

import (
	"context"
	"shop/pkg/querier"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var trKey = struct{}{}

type TrFabricImpl struct {
	querier querier.Querier
}

func NewTrFabric(querier querier.Querier) Fabric {
	return &TrFabricImpl{
		querier: querier,
	}
}

func (tf *TrFabricImpl) Begin(ctx context.Context) (context.Context, Transaction, error) {
	tx, err := tf.querier.BeginRepeatableTx(ctx)
	if err != nil {
		return nil, nil, err
	}
	tr := newTransaction(tx)
	return setContextTr(ctx, tr), tr, nil
}

var _ querier.Querier = &transaction{}

type transaction struct {
	tx pgx.Tx
}

func newTransaction(tx pgx.Tx) Transaction {
	return &transaction{
		tx: tx,
	}
}

func (t *transaction) BeginRepeatableTx(ctx context.Context) (pgx.Tx, error) {
	return t.tx, nil
}

func (t *transaction) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, sql, arguments...)
}

func (t *transaction) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *transaction) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func setContextTr(ctx context.Context, tr Transaction) context.Context {
	return context.WithValue(ctx, trKey, tr)
}

func Get(ctx context.Context, qr querier.Querier) querier.Querier {
	tr := ctx.Value(trKey)
	if tr == nil {
		return qr
	}
	return tr.(*transaction)
}
