package user

import (
	"context"
	"fmt"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/repository/common"
	scheme_user "shop/internal/repository/scheme/user"
	"shop/pkg/querier"
	"shop/pkg/transaction"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrNoToken      = errors.New("no token provided")
    ErrExpiredToken = errors.New("token expired")
	ErrNotFound = errors.New("user not auth")
)

type UserClaims struct {
    jwt.RegisteredClaims
    UserID uuid.UUID `json:"user_id"`
}

type Repository struct {
	querier querier.Querier
	secret   *conf.Secrets
}


func NewRepository(querier querier.Querier, secret *conf.Secrets) *Repository {
	return &Repository{
		querier: querier,
		secret: secret,
	}
}

func (repo *Repository) Get(ctx context.Context, filter Filter, opts GetOptions) (models.User, error) {
	query := sq.Select(scheme_user.ID, scheme_user.Name, scheme_user.Password, scheme_user.Balance).
		PlaceholderFormat(sq.Dollar).
		From(scheme_user.Table)

	if filter.Username != nil {
		query = query.Where(sq.Eq{scheme_user.Name: *filter.Username})
	}

	if filter.ID != nil {
		query = query.Where(sq.Eq{scheme_user.ID: *filter.ID})
	}

	if opts.ForUpdate {
		query = query.Suffix("FOR UPDATE")
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return models.User{}, fmt.Errorf("building sql: %w", err)
	}

	row := transaction.Get(ctx, repo.querier).QueryRow(ctx, sql, args...)
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

func (repo *Repository) Update(ctx context.Context, update Update, filter Filter) error {
	query := sq.Update(scheme_user.Table).PlaceholderFormat(sq.Dollar)

	if filter.ID != nil {
		query = query.Where(sq.Eq{scheme_user.ID: filter.ID})
	}

	if filter.Username != nil {
		query = query.Where(sq.Eq{scheme_user.Name: filter.Username})
	}

	if update.Balance != nil {
		query = query.Set(scheme_user.Balance, update.Balance)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building sql: %w", err)
	}
	
	_, err = transaction.Get(ctx, repo.querier).Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("executing: %w", err)
	}

	return nil
} 

// extractToken извлекает токен из заголовка Authorization
func extractToken(authorization string) (string, error) {
    if authorization == "" {
        return "", ErrNoToken
    }

    parts := strings.Split(authorization, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", ErrInvalidToken
    }

    return parts[1], nil
}

// IsAuth проверяет валидность токена 
func (repo *Repository) IsAuth(ctx context.Context) (uuid.UUID, error) {
    tc, ok := transport.FromServerContext(ctx)
	if !ok {
        return uuid.Nil, ErrNoToken
    }
	headers := tc.RequestHeader()
    authorization := headers.Get("authorization")
    tokenString, err := extractToken(authorization)
    if err != nil {
        return uuid.Nil, err
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return []byte(repo.secret.JwtKey), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return uuid.Nil, ErrExpiredToken
        }
        return uuid.Nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return uuid.Nil, ErrInvalidToken
    }

    // Проверяем наличие необходимых полей
    userIDStr, ok := claims["id"].(string)
    if !ok {
        return uuid.Nil, ErrInvalidToken
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, ErrInvalidToken
    }

    return userID, nil
}
