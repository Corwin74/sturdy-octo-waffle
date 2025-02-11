package querier

import (
	"fmt"


	"github.com/jackc/pgx"
)

type PostgresConf struct {
	Source string
}


func NewPostgres(conf PostgresConf) (Querier, error) {
	config, err := pgx.ParseDSN(conf.Source)
	if err != nil {
		return nil, fmt.Errorf("parsing DSN database: %w", err)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("unable connect to database: %w", err)
	}

	return conn, nil
}
