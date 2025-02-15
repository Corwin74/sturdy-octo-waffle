package transaction

import (
	"context"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Fabric interface {
	Begin(ctx context.Context) (context.Context, Transaction, error)
}
