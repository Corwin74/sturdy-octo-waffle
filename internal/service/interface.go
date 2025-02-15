package service

import (
	"context"
	)

type UserUsecase interface {
	Auth(ctx context.Context, username, password string) (string, error)
	TransferCoins(ctx context.Context, toUser string, amount uint) error  
}
