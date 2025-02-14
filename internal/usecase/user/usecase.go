package user

import (
	"context"
	"errors"
	"fmt"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/repository/common"
	repo_user "shop/internal/repository/user"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// Usecase -- пользователя
type Usecase struct {
	userRepo UserRepo
	config   *conf.Secrets
}

// NewUsecase -- конструктор
func NewUsecase(userRepo UserRepo, config *conf.Secrets) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		config:   config,
	}
}

// Auth -- авторизует пользователя. если пользователь не найден то создает его
//
// возвращет JWT
func (uc *Usecase) Auth(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.Get(ctx, repo_user.Filter{Username: &username})

	userID := user.ID

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			hashedPassword, err := HashPassword(password)
			if err != nil {
				return "", fmt.Errorf("hashing password: %w", err)
			}
			user := models.User{
				Name:     username,
				Password: hashedPassword,
				Balance:  1000,
			}
			id, err := uc.userRepo.Create(ctx, user)
			if err != nil {
				return "", fmt.Errorf("creating user: %w", err)
			}
			userID = id
		} else {
			return "", fmt.Errorf("getting user: %w", err)
		}
	} else {
		if !checkPassword(password, user.Password) {
			return "", fmt.Errorf("invalid password")
		}
	}

	return generateTokenForUser(userID, uc.config.JwtKey), nil
}

func generateTokenForUser(userID uuid.UUID, secret string) string {
	now := time.Now()

	claims := jwt.MapClaims{
		"id":  userID.String(),
		"iat": now.Unix(),                     // время создания токена
		"exp": now.Add(24 * time.Hour).Unix(), // срок действия - 24 часа
		"jti": uuid.New().String(),            // уникальный идентификатор токена
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

// HashPassword -- хеширует пароль для безопасного хранения
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
