package querier

import (
	"context"
	"fmt"
	"shop/internal/conf"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type PostgresConf struct {
// 	Source string
// }

// Database - структура, реализующая интерфейс Querier
type Database struct {
    pool *pgxpool.Pool
}

// func NewPostgres(conf PostgresConf) (Querier, error) {
// 	config, err := pgx.ParseDSN(conf.Source)
// 	if err != nil {
// 		return nil, fmt.Errorf("parsing DSN database: %w", err)
// 	}

// 	conn, err := pgx.Connect(config)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable connect to database: %w", err)
// 	}

// 	return conn, nil
// }


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