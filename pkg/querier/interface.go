package querier

import (
	"context"

	"github.com/jackc/pgx/v5"
)


type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

