package querier

import "github.com/jackc/pgx"


type Querier interface {
	QueryRow(sql string, args ...interface{}) *pgx.Row
}

