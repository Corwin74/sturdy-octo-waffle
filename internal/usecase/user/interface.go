package user

import (
	"context"
	"shop/internal/models"
	repo_user "shop/internal/repository/user"
	repo_item "shop/internal/repository/item"
	repo_useritem "shop/internal/repository/useritem"
	repo_transferhistory "shop/internal/repository/transferhistory"

	"github.com/google/uuid"
)

// UserRepo -- репозиторий для работы с User
type UserRepo interface {
	Get(ctx context.Context, filter repo_user.Filter, opts repo_user.GetOptions) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
	Update(ctx context.Context, update repo_user.Update, filter repo_user.Filter) error
	IsAuth(ctx context.Context) (uuid.UUID, error)
}

type ItemRepo interface {
	Get(ctx context.Context, filter repo_item.Filter) (models.Item, error)
	GetMany(ctx context.Context, filter repo_item.Filter) ([]models.Item, error)
}

type TransferHistory interface{
	Create(ctx context.Context, th models.TransferHistory) (uuid.UUID, error)
	GetMany(ctx context.Context, filter repo_transferhistory.Filter) ([]models.TransferHistory, error)
}

type UserItemRepo interface {
	Get(ctx context.Context, filter repo_useritem.Filter) (models.UserItem, error)
	Create(ctx context.Context, md models.UserItem) (uuid.UUID, error)
	GetUserItemsAmount(ctx context.Context, filter repo_useritem.Filter) ([]models.UserItemsAmount, error)
}
