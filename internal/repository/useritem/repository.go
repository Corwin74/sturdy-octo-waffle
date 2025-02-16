package useritem

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/repository/common"
	"shop/internal/repository/scheme/useritem"
	"shop/pkg/querier"
	"shop/pkg/transaction"

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

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.UserItem, error) {
	query := sq.Select(
		scheme_useritem.ID,
		scheme_useritem.UserID,
		scheme_useritem.ItemID,
	).From(scheme_useritem.Table)

	if filter.UserID != nil {
		query = query.Where(sq.Eq{scheme_useritem.UserID: *filter.UserID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.UserItem{}, fmt.Errorf("building sql: %w", err)
	}

	row := transaction.Get(ctx, repo.querier).QueryRow(ctx, sql, args...)
	var dbModel scheme_useritem.UserItem
	err = row.Scan(&dbModel.ID, &dbModel.UserID, &dbModel.ItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserItem{}, common.ErrNotFound
		}

		return models.UserItem{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return models.UserItem{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}

func (repo *Repository) Create(ctx context.Context, md models.UserItem) (uuid.UUID, error) {
	dbModel := scheme_useritem.ConvertToDBModel(md)

	query := sq.Insert(scheme_useritem.Table).PlaceholderFormat(sq.Dollar).
	Columns(
		scheme_useritem.UserID,
		scheme_useritem.ItemID,
	).Values(
		dbModel.UserID,
		dbModel.ItemID,
	).Suffix("RETURNING " + scheme_useritem.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("building sql: %w", err)
	}

	row := transaction.Get(ctx, repo.querier).QueryRow(ctx, sql, args...)
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

func (repo *Repository) GetUserItemsAmount(ctx context.Context, filter Filter) ([]models.UserItemsAmount, error) {
	query := sq.Select(
		scheme_useritem.ItemID,
		"count(" + scheme_useritem.ItemID + ")",
	).PlaceholderFormat(sq.Dollar).From(scheme_useritem.Table)
	
	if filter.UserID != nil {
		query = query.Where(sq.Eq{scheme_useritem.UserID: *filter.UserID})
	}
	
	query = query.GroupBy(scheme_useritem.ItemID)
	
	sql, args, err := query.ToSql()
	if err != nil {
		return []models.UserItemsAmount{}, fmt.Errorf("building sql: %w", err)
	}

	rows, err := repo.querier.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.UserItemsAmount{}, common.ErrNotFound
		}

		return []models.UserItemsAmount{}, fmt.Errorf("quering: %w", err)
	}
	var domainModels []models.UserItemsAmount
	defer rows.Close()

	for rows.Next() {
		var dbModel scheme_useritem.UserItemsAmount
		err = rows.Scan(&dbModel.ItemID, &dbModel.Quantity)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []models.UserItemsAmount{}, common.ErrNotFound
			}
			return []models.UserItemsAmount{}, fmt.Errorf("scanning: %w", err)
		}
		domainModel, _ := dbModel.ConvertToDomainModel()
		domainModels = append(domainModels, domainModel)
	}

	return domainModels, nil
}
