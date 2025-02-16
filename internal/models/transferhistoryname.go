package models

import "github.com/google/uuid"


type TransferHistoryName struct {
	ID uuid.UUID
	SenderName string
	ReceiverName string
	Amount uint
}
