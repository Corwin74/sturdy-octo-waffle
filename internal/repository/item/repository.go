package item

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/repository/common"
	scheme_item "shop/internal/repository/scheme/item"
	"shop/pkg/querier"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)


type Repository struct {
	querier querier.Querier
}



func NewRepository(querier querier.Querier) *Repository {
	return &Repository{
		querier: querier,
	}
}

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.Item, error) {
	query := sq.Select(scheme_item.ID, scheme_item.Name, scheme_item.Price).
		From(scheme_item.Table)

	if filter.Name != nil {
		query = query.Where(sq.Eq{scheme_item.Name: *filter.Name})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Item{}, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var dbModel scheme_item.Item
	err = row.Scan(&dbModel.ID, &dbModel.Name, &dbModel.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Item{}, common.ErrNotFound
		}

		return models.Item{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return  models.Item{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}


func (repo *Repository) Create(ctx context.Context, user models.Item) (uuid.UUID, error) {
	dbModel := scheme_item.ConvertToDBModel(user)

	query := sq.Insert(scheme_item.Table).Columns(
		scheme_item.Name,
		scheme_item.Price,
	).Values(
		dbModel.Name,
		dbModel.Price,
	).Suffix("RETURNING " + scheme_item.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var idString string
	err = row.Scan(&idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("quering: %w", err)
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing id: %w", err)
	}
	return id, nil
}
