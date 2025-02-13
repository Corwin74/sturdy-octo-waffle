package user

import (
	"context"
	"errors"
	"fmt"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/repository/common"
	repo_user "shop/internal/repository/user"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type Usecase struct {
	userRepo UserRepo
	config *conf.Secrets
}


func NewUsecase(userRepo UserRepo, config *conf.Secrets) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		config: config,
	}
}


// Auth -- авторизует пользователя. если пользователь не найден то создает его
//
// возвращет идентификатор юзера
func (u *Usecase) Auth(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.Get(ctx, repo_user.Filter{Username: &username})
	userID := user.ID

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			user := models.User{
				Name: username,
				Password: password,
				Balance: 1000,
			}
			id, err := u.userRepo.Create(ctx, user)
			if err != nil {
				return "", fmt.Errorf("creating user: %w", err)
			}
			userID = id
		} else {
			return "", fmt.Errorf("getting user: %w", err)
		}

		
	}
	
	token, err := generateTokenForUser(userID, u.config.JwtKey)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	return token, nil
}

func generateTokenForUser(userID uuid.UUID, secret string) (string, error) {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": userID.String(),
		}).SignedString([]byte(secret))
	return token, err
}
