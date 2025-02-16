package transferhistoryname

import (
	"context"
	"fmt"
	"shop/internal/models"
	"shop/internal/common"
	scheme_transferhistoryname "shop/internal/repository/scheme/transferhistoryname"
	scheme_user "shop/internal/repository/scheme/user"
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

func (repo *Repository) Get(ctx context.Context, filter Filter) (models.TransferHistoryName, error) {
	query := sq.Select(
		scheme_transferhistoryname.ID,
		scheme_transferhistoryname.SenderName,
		scheme_transferhistoryname.ReceiverName,
		scheme_transferhistoryname.Amount,
	).PlaceholderFormat(sq.Dollar).From(scheme_user.Table)

	if filter.SenderName != nil {
		query = query.Where(sq.Eq{scheme_transferhistoryname.SenderName: *filter.SenderName})
	}

	if filter.ReceiverName != nil {
		query = query.Where(sq.Eq{scheme_transferhistoryname.ReceiverName: *filter.ReceiverName})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return models.TransferHistoryName{}, fmt.Errorf("building sql: %w", err)
	}

	row := repo.querier.QueryRow(ctx, sql, args...)
	var dbModel scheme_transferhistoryname.TransferHistoryName
	err = row.Scan(&dbModel.ID, &dbModel.SenderName, &dbModel.ReceiverName, &dbModel.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TransferHistoryName{}, common.ErrNotFound
		}

		return models.TransferHistoryName{}, fmt.Errorf("quering: %w", err)
	}

	domainModel, err := dbModel.ConvertToDomainModel()
	if err != nil {
		return models.TransferHistoryName{}, fmt.Errorf("converting: %w", err)
	}

	return domainModel, nil
}

func (repo *Repository) GetMany(ctx context.Context, filter Filter) ([]models.TransferHistoryName, error) {
	query := sq.Select(
		scheme_transferhistoryname.ID,
		scheme_transferhistoryname.SenderName,
		scheme_transferhistoryname.ReceiverName,
		scheme_transferhistoryname.Amount,
	).PlaceholderFormat(sq.Dollar).From(scheme_transferhistoryname.Table)

	if filter.SenderName != nil {
		query = query.Where(sq.Eq{scheme_transferhistoryname.SenderName: *filter.SenderName})
	}

	if filter.ReceiverName != nil {
		query = query.Where(sq.Eq{scheme_transferhistoryname.ReceiverName: *filter.ReceiverName})
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
	var dbModels []models.TransferHistoryName

	for rows.Next() {
		var dbModel scheme_transferhistoryname.TransferHistoryName
		err = rows.Scan(&dbModel.ID, &dbModel.SenderName, &dbModel.ReceiverName, &dbModel.Amount)
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

func (repo *Repository) Create(ctx context.Context, th models.TransferHistoryName) (uuid.UUID, error) {
	dbModel := scheme_transferhistoryname.ConvertToDBModel(th)

	query := sq.Insert(scheme_transferhistoryname.Table).PlaceholderFormat(sq.Dollar).
		Columns(
			scheme_transferhistoryname.SenderName,
			scheme_transferhistoryname.ReceiverName,
			scheme_transferhistoryname.Amount,
		).Values(
		dbModel.SenderName,
		dbModel.ReceiverName,
		dbModel.Amount,
	).Suffix("RETURNING " + scheme_transferhistoryname.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("building sql: %w", err)
	}

	row := transaction.Get(ctx, repo.querier).QueryRow(ctx, sql, args...)
	// row := repo.querier.QueryRow(ctx, sql, args...)
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
