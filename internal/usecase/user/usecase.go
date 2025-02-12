package user

import (
	"context"

)


type Usecase struct {
	userRepo UserRepo
}


func NewUsecase(userRepo UserRepo) *Usecase {
	return &Usecase{
		userRepo: userRepo,
	}
}

func (u *Usecase) Auth(ctx context.Context) error {
	return nil
}
