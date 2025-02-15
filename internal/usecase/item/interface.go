package item

import (
	"context"
	"shop/internal/models"
	repo_item "shop/internal/repository/item"
)

// ItemRepo -- репозиторий мерча
type ItemRepo interface {
	Get(ctx context.Context, filter repo_item.Filter, opts repo_item.GetOptions) (models.Item, error)
}
