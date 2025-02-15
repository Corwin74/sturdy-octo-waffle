package service

import (
	"context"
	"shop/internal/models"
	)

type UserUsecase interface {
	Auth(ctx context.Context, username, password string) (string, error)
	IsAuth(ctx context.Context) (models.User, error)
	TransferCoins(ctx context.Context, toUser string, amount uint) error  
}
