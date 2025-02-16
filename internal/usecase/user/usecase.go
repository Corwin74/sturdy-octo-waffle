package user

import (
	"context"
	"errors"
	"fmt"
	"shop/internal/conf"
	"shop/internal/models"
	"shop/internal/common"
	repo_item "shop/internal/repository/item"
	repo_user "shop/internal/repository/user"
	repo_useritem "shop/internal/repository/useritem"
	repo_transferhistory "shop/internal/repository/transferhistoryname"
	"shop/pkg/querier"
	"shop/pkg/transaction"

	"github.com/google/uuid"
)

// Usecase -- пользователя
type Usecase struct {
	userRepo          UserRepo
	itemRepo 		  ItemRepo
	userItemRepo      UserItemRepo
	transferHistoryName   TransferHistoryName
	config            *conf.Secrets
	querier           querier.Querier
	transactionFabric transaction.Fabric
}

// NewUsecase -- конструктор
func NewUsecase(
	userRepo UserRepo,
	transferHistoryName TransferHistoryName,
	itemRepo ItemRepo,
	userItemRepo UserItemRepo,
	config *conf.Secrets,
	querier querier.Querier,
	transactionFabric transaction.Fabric,
) *Usecase {
	return &Usecase{
		userRepo:        userRepo,
		itemRepo: 		 itemRepo,
		userItemRepo:    userItemRepo,	
		transferHistoryName: transferHistoryName,
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
			return "", common.ErrUnauthorized
		}
	}

	return generateTokenForUser(userID, uc.config.JwtKey), nil
}

func (uc *Usecase) TransferCoins(ctx context.Context, toUserName string, amount uint) error {
	fromUserID, err := uc.userRepo.IsAuth(ctx)
	if err != nil {
		return common.ErrUnauthorized
	}

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
	transferHistory := models.TransferHistoryName{
		SenderName:   fromUser.Name,
		ReceiverName: toUser.Name,
		Amount:     amount,
	}
	_, err = uc.transferHistoryName.Create(ctx, transferHistory)
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

func (uc *Usecase)Buy(ctx context.Context, itemName string) error {
	userID, err := uc.userRepo.IsAuth(ctx)
	if err != nil {
		return common.ErrUnauthorized
	}

	ctx, tr, err := uc.transactionFabric.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	userFilter := repo_user.Filter{ID: &userID}
	itemFilter := repo_item.Filter{Name: &itemName}
	getOpts := repo_user.GetOptions{ForUpdate: true}
	user, err := uc.userRepo.Get(ctx, userFilter, getOpts)
	if err != nil {
		return fmt.Errorf("receive user: %w", err)
	}
	item, err := uc.itemRepo.Get(ctx, itemFilter)
	if err != nil {
		return fmt.Errorf("item: %w", err)
	}
	if user.Balance < item.Price {
		return fmt.Errorf("no enough money: %w", err)
	}

	user.Balance -= item.Price

	err = uc.userRepo.Update(ctx, repo_user.Update{Balance: &user.Balance}, userFilter)
	if err != nil {
		if tErr := tr.Rollback(ctx); tErr != nil {
			return fmt.Errorf("rollbacking transaction (%s): %w", err, tErr)
		}
		return fmt.Errorf("commiting transaction: %w", err)
	}
	_, err = uc.userItemRepo.Create(ctx, models.UserItem{UserID: userID, ItemID: item.ID})
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

func (uc *Usecase)Info(ctx context.Context) (models.UserInfo, error) {
	userID, err := uc.userRepo.IsAuth(ctx)
	if err != nil {
		return models.UserInfo{}, common.ErrUnauthorized
	}

	ctx, tr, err := uc.transactionFabric.Begin(ctx)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("begin transaction: %w", err)
	}

	userFilter := repo_useritem.Filter{UserID: &userID}
	getOpts := repo_user.GetOptions{ForUpdate: true}
	user, err := uc.userRepo.Get(ctx, repo_user.Filter{ID: &userID}, getOpts)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("info user: %w", err)
	}
	
	userItemsAmount, err := uc.userItemRepo.GetUserItemsAmount(ctx, userFilter)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("get user item amount: %w", err)
	}

	items, err := uc.itemRepo.GetMany(ctx, repo_item.Filter{})
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("get items: %w", err)
	}

	idToName := make(map[uuid.UUID]string)

	for _, item := range items {
		idToName[item.ID] = item.Name
	}

	var inventory []models.InventoryItem
	for _, itemAmount := range userItemsAmount {
		name, exist := idToName[itemAmount.ItemID]
		if !exist {
			return models.UserInfo{}, fmt.Errorf("item name not found: %w", err)
		}
		inventory = append(inventory, models.InventoryItem{Type: name, Quantity: uint(itemAmount.Quantity)})
	}

	receivedCoinsHistory, err := uc.transferHistoryName.
		GetMany(ctx, repo_transferhistory.Filter{ReceiverName: &user.Name})
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("receive history: %w", err)
	}

	sentCoinsHistory, err := uc.transferHistoryName.
	GetMany(ctx, repo_transferhistory.Filter{SenderName: &user.Name})
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("send history: %w", err)
	}
	err = tr.Commit(ctx)
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("commiting transaction: %w", err)
	}
	var sentCoins []models.SentCoins
	for _, transfer := range sentCoinsHistory {
		sentCoins = append(sentCoins, models.SentCoins{ToUser: transfer.ReceiverName, Amount: transfer.Amount})
	}
	
	var receivedCoins []models.ReceivedCoins
	for _, transfer := range receivedCoinsHistory {
		receivedCoins = append(receivedCoins, models.ReceivedCoins{FromUser: transfer.SenderName, Amount: transfer.Amount})
	}
	
	coinHistory := models.CoinHistory{Received: receivedCoins, Sent: sentCoins}

	return models.UserInfo{
		Coins: user.Balance,
		Inventory: inventory,
		CoinHistory: coinHistory,
	}, nil
}
