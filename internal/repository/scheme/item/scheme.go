package scheme_item

import (
	"fmt"
	"shop/internal/models"

	"github.com/google/uuid"
)

const (
	Table = "items"

	ID = "id"
	Name = "name"
	Price = "price"
)


type Item struct {
	ID string
	Name string
	Price int
}


func ConvertToDBModel(item models.Item) Item {
	return Item{
		ID: item.ID.String(),
		Name: item.Name,
		Price: int(item.Price),
	}
}


func (i Item) ConvertToDomainModel() (models.Item, error) {
	id, err := uuid.Parse(i.ID)
	if err != nil {
		return models.Item{}, fmt.Errorf("parsing id: %w", err)
	}

	return models.Item{
		ID: id,
		Name: i.Name,
		Price: uint(i.Price),
	}, nil
}
