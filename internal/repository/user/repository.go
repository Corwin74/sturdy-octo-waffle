package user

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/repository/common"
	scheme_user "shop/internal/repository/scheme/user"
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

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.User, error) {
	query := sq.Select(scheme_user.ID, scheme_user.Name, scheme_user.Password, scheme_user.Balance).
		PlaceholderFormat(sq.Dollar).
		From(scheme_user.Table)

	if filter.Username != nil {
		query = query.Where(sq.Eq{scheme_user.Name: *filter.Username})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.User{}, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var dbModel scheme_user.User
	err = row.Scan(&dbModel.ID, &dbModel.Name, &dbModel.Password, &dbModel.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, common.ErrNotFound
		}

		return models.User{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return  models.User{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}


func (repo *Repository) Create(ctx context.Context, user models.User) (uuid.UUID, error) {
	dbModel := scheme_user.ConvertToDBModel(user)
	query := sq.Insert(scheme_user.Table).PlaceholderFormat(sq.Dollar).Columns(
		scheme_user.Name,
		scheme_user.Password,
		scheme_user.Balance,
	).Values(
		dbModel.Name,
		dbModel.Password,
		dbModel.Balance,
	).Suffix("RETURNING " + scheme_user.ID)

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
