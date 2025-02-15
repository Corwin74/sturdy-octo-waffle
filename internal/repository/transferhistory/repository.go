package transferhistory

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/repository/common"
	scheme_transferhistory "shop/internal/repository/scheme/transferhistory"
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

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.TransferHistory, error) {
	query := sq.Select(
		scheme_transferhistory.ID,
		scheme_transferhistory.SenderID,
		scheme_transferhistory.ReceiverID,
		scheme_transferhistory.Amount,
	).PlaceholderFormat(sq.Dollar).From(scheme_user.Table)

	if filter.SenderID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.SenderID: *filter.SenderID})
	}

	if filter.ReceiverID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.ReceiverID: *filter.ReceiverID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var dbModel scheme_transferhistory.TransferHistory
	err = row.Scan(&dbModel.ID, &dbModel.SenderID, &dbModel.ReceiverID, &dbModel.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TransferHistory{}, common.ErrNotFound
		}

		return models.TransferHistory{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}

func (repo *Repository) GetMany(ctx context.Context, filter Filter) ([]models.TransferHistory, error) {
	query := sq.Select(
		scheme_transferhistory.ID,
		scheme_transferhistory.SenderID,
		scheme_transferhistory.ReceiverID,
		scheme_transferhistory.Amount,
	).PlaceholderFormat(sq.Dollar).From(scheme_user.Table)

	if filter.SenderID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.SenderID: *filter.SenderID})
	}

	if filter.ReceiverID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.ReceiverID: *filter.ReceiverID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building sql: %w", err)
	}

	rows, err := repo.querier.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("quering: %w", err)
	}
	defer rows.Close()
	var dbModels []models.TransferHistory

	for rows.Next() {
		var dbModel scheme_transferhistory.TransferHistory
		err = rows.Scan(&dbModel.ID, &dbModel.SenderID, &dbModel.ReceiverID, &dbModel.Amount)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, common.ErrNotFound
			}
		}

		domainModel, err := dbModel.ConvertToDomainModel()
		if err != nil {
			return nil, fmt.Errorf("converting: %w", err)
		}
		dbModels = append(dbModels, domainModel)

	}

	return dbModels, nil
}

func (repo *Repository) Create(ctx context.Context, th models.TransferHistory) (uuid.UUID, error) {
	dbModel := scheme_transferhistory.ConvertToDBModel(th)

	query := sq.Insert(scheme_transferhistory.Table).PlaceholderFormat(sq.Dollar).
		Columns(
			scheme_transferhistory.SenderID,
			scheme_transferhistory.ReceiverID,
			scheme_transferhistory.Amount,
		).Values(
		dbModel.SenderID,
		dbModel.ReceiverID,
		dbModel.Amount,
	).Suffix("RETURNING " + scheme_transferhistory.ID)

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

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.TransferHistory, error) {
	query := sq.Select(
		scheme_transferhistory.ID,
		scheme_transferhistory.SenderID,
		scheme_transferhistory.ReceiverID,
		scheme_transferhistory.Amount,
	).PlaceholderFormat(sq.Dollar).From(scheme_user.Table)

	if filter.SenderID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.SenderID: *filter.SenderID})
	}

	if filter.ReceiverID != nil {
		query = query.Where(sq.Eq{scheme_transferhistory.ReceiverID: *filter.ReceiverID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var dbModel scheme_transferhistory.TransferHistory
	err = row.Scan(&dbModel.ID, &dbModel.SenderID, &dbModel.ReceiverID, &dbModel.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TransferHistory{}, common.ErrNotFound
		}

		return models.TransferHistory{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}
