package models

import "github.com/google/uuid"


type UserItem struct {
	ID uuid.UUID
	UserID uuid.UUID
	ItemID uuid.UUID
}