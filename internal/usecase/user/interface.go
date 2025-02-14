package user

import (
	"context"
	"shop/internal/models"
	repo_user "shop/internal/repository/user"

	"github.com/google/uuid"
)

// UserRepo -- репозиторий для работы с User
type UserRepo interface {
	Get(ctx context.Context, filter repo_user.Filter) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
}
