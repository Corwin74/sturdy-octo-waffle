package querier

import (
	"context"
	"fmt"
	"shop/internal/conf"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database - структура, реализующая интерфейс Querier
type Database struct {
    pool *pgxpool.Pool
}

// NewDatabase создает новый экземпляр Database
func NewDatabase(conf *conf.Data) (*Database, error) {
    ctx := context.Background()
    
	if conf == nil {
		return nil, fmt.Errorf("no config data")
	}
	
	config, err := pgxpool.ParseConfig(conf.Database.Source)
	if err != nil {
		return nil, fmt.Errorf("parsing DSN database: %w", err)
	}

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, err
    }

    return &Database{
        pool: pool,
    }, nil
}


// QueryRow реализует метод интерфейса Querier
func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
    return db.pool.QueryRow(ctx, sql, args...)
}

// Query реализует метод интерфейса Querier
func (db *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
    return db.pool.Query(ctx, sql, args...)
}

// Close закрывает соединение с базой данных
func (db *Database) Close() {
    db.pool.Close()
}

// BeginRepeatableTx - начинает транзакцию
func (db *Database) BeginRepeatableTx(ctx context.Context) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
}

func (db *Database) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

var _ Querier = (*Database)(nil)
