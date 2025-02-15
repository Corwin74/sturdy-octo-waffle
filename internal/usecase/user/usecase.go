package user

import (
	"context"
	"errors"
	"fmt"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/repository/common"
	repo_user "shop/internal/repository/user"
	"shop/pkg/querier"
	"shop/pkg/transaction"
)

// Usecase -- пользователя
type Usecase struct {
	userRepo          UserRepo
	transferHistory   TransferHistory
	config            *conf.Secrets
	querier           querier.Querier
	transactionFabric transaction.Fabric
}

// NewUsecase -- конструктор
func NewUsecase(
	userRepo UserRepo,
	transferHistory TransferHistory,
	config *conf.Secrets,
	querier querier.Querier,
	transactionFabric transaction.Fabric,
) *Usecase {
	return &Usecase{
		userRepo:        userRepo,
		transferHistory: transferHistory,
		config:          config,
		querier:         querier,
		transactionFabric: transactionFabric,
	}
}

// Auth -- авторизует пользователя. если пользователь не найден то создает его
//
// возвращет JWT
func (uc *Usecase) Auth(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.Get(ctx, repo_user.Filter{Username: &username}, repo_user.GetOptions{})

	userID := user.ID

	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			hashedPassword, err := GenerateArgon2Hash(password)
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
		if match, _ := verifyPassword(password, user.Password); !match {
			return "", fmt.Errorf("invalid password")
		}
	}

	return generateTokenForUser(userID, uc.config.JwtKey), nil
}

func (uc *Usecase) TransferCoins(ctx context.Context, toUserName string, amount uint) error {
	fromUserID, err := uc.userRepo.IsAuth(ctx)
	ctx, tr, err := uc.transactionFabric.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	toUserFilter := repo_user.Filter{Username: &toUserName}
	fromUserFilter := repo_user.Filter{ID: &fromUserID}
	getOpts := repo_user.GetOptions{ForUpdate: true}
	toUser, err := uc.userRepo.Get(ctx, toUserFilter, getOpts)
	if err != nil {
		return fmt.Errorf("receive user: %w", err)
	}
	fromUser, err := uc.userRepo.Get(ctx, fromUserFilter, getOpts)
	if err != nil {
		return fmt.Errorf("receive user: %w", err)
	}
	if fromUser.Balance < amount {
		return fmt.Errorf("no enough money: %w", err)
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	err = uc.userRepo.Update(ctx, repo_user.Update{Balance: &fromUser.Balance}, fromUserFilter)
	if err != nil {
		if tErr := tr.Rollback(ctx); tErr != nil {
			return fmt.Errorf("rollbacking transaction (%s): %w", err, tErr)
		}
		return fmt.Errorf("commiting transaction: %w", err)
	}

	err = uc.userRepo.Update(ctx, repo_user.Update{Balance: &toUser.Balance}, toUserFilter)
	if err != nil {
		if tErr := tr.Rollback(ctx); tErr != nil {
			return fmt.Errorf("rollbacking transaction (%s): %w", err, tErr)
		}
		return fmt.Errorf("commiting transaction: %w", err)
	}
	transferHistoty := models.TransferHistory{
		SenderID:   fromUser.ID,
		ReceiverID: toUser.ID,
		Amount:     amount,
	}
	_, err = uc.transferHistory.Create(ctx, transferHistoty)
	if err != nil {
		if tErr := tr.Rollback(ctx); tErr != nil {
			return fmt.Errorf("rollbacking transaction (%s): %w", err, tErr)
		}
		return fmt.Errorf("commiting transaction: %w", err)
	}
	err = tr.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}
	return nil
}
