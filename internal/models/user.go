package models

import "github.com/google/uuid"


type User struct {
	ID uuid.UUID
	Name string
	Password string
	Balance uint
}

type UserInfo struct {
    Coins       uint             `json:"coins"`
    Inventory   []InventoryItem  `json:"inventory"`
    CoinHistory CoinHistory      `json:"coinHistory"`
}

type InventoryItem struct {
    Type     string `json:"type"`
    Quantity uint   `json:"quantity"`
}

type CoinHistory struct {
    Received []ReceivedCoins `json:"received"`
    Sent     []SentCoins     `json:"sent"`
}

type ReceivedCoins struct {
    FromUser string `json:"fromUser"`
    Amount   uint   `json:"amount"`
}

type SentCoins struct {
    ToUser string `json:"toUser"`
    Amount uint   `json:"amount"`
}
