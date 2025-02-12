package user

import (
	"shop/internal/usecase/user"
	"shop/pkg/querier"
)


type Repository struct {
	querier querier.Querier
}

var _ user.UserRepo = (*Repository)(nil)


func NewRepository(querier querier.Querier) *Repository {
	return &Repository{
		querier: querier,
	}
}

