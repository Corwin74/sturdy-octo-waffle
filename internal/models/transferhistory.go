package models

import "github.com/google/uuid"


type TransferHistory struct {
	ID uuid.UUID
	SenderID uuid.UUID
	ReceiverID uuid.UUID
	Amount uint
}
