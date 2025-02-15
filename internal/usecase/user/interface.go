package user

import (
	"context"
	"shop/internal/models"
	repo_user "shop/internal/repository/user"

	"github.com/google/uuid"
)

// UserRepo -- репозиторий для работы с User
type UserRepo interface {
	Get(ctx context.Context, filter repo_user.Filter, opts repo_user.GetOptions) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
	Update(ctx context.Context, update repo_user.Update, filter repo_user.Filter) error
	IsAuth(ctx context.Context) (models.User, error)
}

type TransferHistory interface{
	Create(ctx context.Context, th models.TransferHistory) (uuid.UUID, error)
}

func UserRepoWithTx(repo UserRepo) {
	
}
