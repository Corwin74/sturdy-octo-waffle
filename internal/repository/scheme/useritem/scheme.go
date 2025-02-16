// CREATE TABLE IF NOT EXISTS users_items (
//     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     id_user UUID NOT NULL REFERENCES users(id) ,
//     id_item UUID NOT NULL REFERENCES items(id)  
// );

package scheme_useritem

import (
	"fmt"
	"shop/internal/models"

	"github.com/google/uuid"
)

const (
	Table = "users_items"

	ID = "id"
	UserID = "id_user"
	ItemID = "id_item"
)


type UserItem struct {
	ID string
	UserID string
	ItemID string
}

type UserItemsAmount struct {
	ItemID string
	Quantity int
}

func ConvertToDBModel(useritem models.UserItem) UserItem {
	return UserItem{
		ID: useritem.ID.String(),
		UserID: useritem.UserID.String(),
		ItemID: useritem.ItemID.String(),
	}
}


func (ui UserItem) ConvertToDomainModel() (models.UserItem, error) {
	id, err := uuid.Parse(ui.ID)
	if err != nil {
		return models.UserItem{}, fmt.Errorf("parsing id: %w", err)
	}
	userID, err := uuid.Parse(ui.UserID)
	if err != nil {
		return models.UserItem{}, fmt.Errorf("parsing userID: %w", err)
	}
	itemID, err := uuid.Parse(ui.ItemID)
	if err != nil {
		return models.UserItem{}, fmt.Errorf("parsing itemID: %w", err)
	}
	return models.UserItem{
		ID: id,
		UserID: userID,
		ItemID: itemID,
	}, nil
}

func (uia UserItemsAmount) ConvertToDomainModel() (models.UserItemsAmount, error) {
	itemID, err := uuid.Parse(uia.ItemID)
	if err != nil {
		return models.UserItemsAmount{}, fmt.Errorf("parsing itemID: %w", err)
	}
	return models.UserItemsAmount{
		ItemID: itemID,
		Quantity: uia.Quantity,
	}, nil
}
