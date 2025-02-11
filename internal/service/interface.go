package service

import "context"

type UserUsecase interface {
	Auth(ctx context.Context) error
}