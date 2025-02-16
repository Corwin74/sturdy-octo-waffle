package item

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/common"
	scheme_item "shop/internal/repository/scheme/item"
	"shop/pkg/querier"
	"shop/pkg/transaction"

	sq "github.com/Masterminds/squirrel"
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
		PlaceholderFormat(sq.Dollar).
		From(scheme_item.Table)

	if filter.Name != nil {
		query = query.Where(sq.Eq{scheme_item.Name: *filter.Name})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Item{}, fmt.Errorf("building sql: %w", err)
	}

	row := transaction.Get(ctx, repo.querier).QueryRow(ctx, sql, args...)
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

func (repo *Repository) GetMany(ctx context.Context, filter Filter) ([]models.Item, error) {
	query := sq.Select(scheme_item.ID, scheme_item.Name, scheme_item.Price).
		PlaceholderFormat(sq.Dollar).
		From(scheme_item.Table)

	if filter.Name != nil {
		query = query.Where(sq.Eq{scheme_item.Name: *filter.Name})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []models.Item{}, fmt.Errorf("building sql: %w", err)
	}

	rows, err := transaction.Get(ctx, repo.querier).Query(ctx, sql, args...)
	defer rows.Close()
	
	var domainModels []models.Item
	for rows.Next() {
		var dbModel scheme_item.Item
		err = rows.Scan(&dbModel.ID, &dbModel.Name, &dbModel.Price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []models.Item{}, common.ErrNotFound
			}
			return []models.Item{}, fmt.Errorf("scanning item: %w", err)
		}
		domainModel, _ :=  dbModel.ConvertToDomainModel()
		domainModels = append(domainModels, domainModel)
	}
	

	return domainModels, nil
}
