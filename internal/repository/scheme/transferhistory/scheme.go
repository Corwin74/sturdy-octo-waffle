package scheme_transferhistory

import (
	"fmt"
	"shop/internal/models"

	"github.com/google/uuid"
)

const (
	Table = "transfers_history"

	ID = "id"
	SenderID = "sender_id"
	ReceiverID = "receiver_id"
	Amount = "amount"
)
// CREATE TABLE IF NOT EXISTS transfers_history (
//     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     sender_id UUID NOT NULL REFERENCES users(id),
//     receiver_id UUID NOT NULL REFERENCES users(id) ,
//     amount INT CHECK(amount BETWEEN 0 AND 3000000) NOT NULL,
//     CONSTRAINT different_users CHECK (sender_id != receiver_id)
// );


type TransferHistory struct {
	ID string
	SenderID string
	ReceiverID string
	Amount int
}


func ConvertToDBModel(th models.TransferHistory) TransferHistory {
	return TransferHistory{
		ID: th.ID.String(),
		SenderID: th.SenderID.String(),
		ReceiverID:  th.ReceiverID.String(),
		Amount: int(th.Amount),
	}
}


func (th TransferHistory) ConvertToDomainModel() (models.TransferHistory, error) {
	id, err := uuid.Parse(th.ID)
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("parsing id: %w", err)
	}
	senderID, err := uuid.Parse(th.SenderID)
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("parsing senderID: %w", err)
	}
	receiiverID, err := uuid.Parse(th.ReceiverID)
	if err != nil {
		return models.TransferHistory{}, fmt.Errorf("parsing receiverID: %w", err)
	}
	return models.TransferHistory{
		ID: id,
		SenderID: senderID,
		ReceiverID: receiiverID,
		Amount: uint(th.Amount),
	}, nil
}
