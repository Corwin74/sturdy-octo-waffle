package scheme_transferhistoryname

import (
	"fmt"
	"shop/internal/models"

	"github.com/google/uuid"
)

const (
	Table = "transfers_history_name"

	ID = "id"
	SenderName = "sender_name"
	ReceiverName = "receiver_name"
	Amount = "amount"
)

type TransferHistoryName struct {
	ID string
	SenderName string
	ReceiverName string
	Amount int
}


func ConvertToDBModel(th models.TransferHistoryName) TransferHistoryName {
	return TransferHistoryName{
		ID: th.ID.String(),
		SenderName: th.SenderName,
		ReceiverName:  th.ReceiverName,
		Amount: int(th.Amount),
	}
}


func (th TransferHistoryName) ConvertToDomainModel() (models.TransferHistoryName, error) {
	id, err := uuid.Parse(th.ID)
	if err != nil {
		return models.TransferHistoryName{}, fmt.Errorf("parsing id: %w", err)
	}
	return models.TransferHistoryName{
		ID: id,
		SenderName: th.SenderName,
		ReceiverName: th.ReceiverName,
		Amount: uint(th.Amount),
	}, nil
}
