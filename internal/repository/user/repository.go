package user

import (
	"shop/pkg/querier"
)


type Repository struct {
	querier querier.Querier
}

func NewRepository(querier querier.Querier) *Repository {
	return &Repository{
		querier: querier,
	}
}

