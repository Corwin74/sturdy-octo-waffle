package user

import (
	"context"
	"errors"
	"shop/internal/common"
	"shop/internal/conf"
	"shop/internal/models"
	repo_item "shop/internal/repository/item"
	repo_transferhistory "shop/internal/repository/transferhistoryname"
	repo_user "shop/internal/repository/user"
	repo_useritem "shop/internal/repository/useritem"
	"shop/mocks/pkg/querier"
	"shop/mocks/pkg/transaction"
	"shop/mocks/usecase"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UsecaseSuite struct {
	suite.Suite
	mc                  *gomock.Controller
	self                *Usecase
	userRepo            *usecase.MockUserRepo
	itemRepo            *usecase.MockItemRepo
	userItemRepo        *usecase.MockUserItemRepo
	transferHistoryName *usecase.MockTransferHistoryName
	config              *conf.Secrets
	querier             *querier.MockQuerier
	transactionFabric   *transaction.MockFabric
	mockTx              *transaction.MockTransaction
}

func TestUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UsecaseSuite))
}

func (u *UsecaseSuite) SetupTest() {
	u.mc = gomock.NewController(u.T())
	u.userRepo = usecase.NewMockUserRepo(u.mc)
	u.itemRepo = usecase.NewMockItemRepo(u.mc)
	u.userItemRepo = usecase.NewMockUserItemRepo(u.mc)
	u.transferHistoryName = usecase.NewMockTransferHistoryName(u.mc)
	u.querier = querier.NewMockQuerier(u.mc)
	u.transactionFabric = transaction.NewMockFabric(u.mc)
	u.mockTx = transaction.NewMockTransaction(u.mc)
	u.config = &conf.Secrets{JwtKey: "some-secret"}

	u.self = NewUsecase(u.userRepo, u.transferHistoryName, u.itemRepo, u.userItemRepo, u.config, u.querier, u.transactionFabric)
}

func (u *UsecaseSuite) TestUsecase_Auth() {
	username := "user1"
	password := "pass"
	user := models.User{ID: uuid.New(), Name: username, Password: "$argon2id$v=19$m=65536,t=1,p=4$FtzfXIrOiVyJai2xrNyDOg$YlkTNrUJub+4XSWg5BsjHs/pLYxJSDxojvLQx9uWkds", Balance: 1000}

	u.Run("success", func() {

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &username}, repo_user.GetOptions{}).
			Return(user, nil)
		token, err := u.self.Auth(context.Background(), username, password)
		u.NoError(err)
		u.True(utf8.RuneCount([]byte(token)) > 10)
	})

	u.Run("wrong password", func() {

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &username}, repo_user.GetOptions{}).
			Return(user, nil)
		token, err := u.self.Auth(context.Background(), username, "kdkdkd")
		u.Error(err)
		u.Equal("", token)
	})

	u.Run("create new user", func() {
		newUsername := "newuser"
		newPassword := "newpass"

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &newUsername}, repo_user.GetOptions{}).
			Return(models.User{}, common.ErrNotFound)

		u.userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, user models.User) (uuid.UUID, error) {
				u.Equal(newUsername, user.Name)
				u.Equal(uint(1000), user.Balance)
				u.True(strings.HasPrefix(user.Password, "$argon2id$"))
				return uuid.New(), nil
			})

		token, err := u.self.Auth(context.Background(), newUsername, newPassword)
		u.NoError(err)
		u.True(utf8.RuneCount([]byte(token)) > 10)
	})

	u.Run("unexected database error", func() {
		username := "user1"
		password := "pass"

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &username}, repo_user.GetOptions{}).
			Return(models.User{}, errors.New("unexpected database error"))

		token, err := u.self.Auth(context.Background(), username, password)
		u.Error(err)
		u.Contains(err.Error(), "getting user")
		u.Equal("", token)
	})
}

func (u *UsecaseSuite) TestUsecase_TransferCoins() {
	ctx := context.Background()
	fromUserID := uuid.New()
	fromUsername := "sender"
	toUsername := "receiver"
	initialBalance := uint(1000)
	transferAmount := uint(500)

	u.Run("успешный перевод монет", func() {
		
		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(fromUserID, nil)

		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &toUsername}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Name: toUsername, Balance: initialBalance}, nil)
		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{ID: &fromUserID}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Name: fromUsername, Balance: initialBalance}, nil)

		u.userRepo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(repo_user.Update{}),
			repo_user.Filter{ID: &fromUserID}).Return(nil)
		u.userRepo.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(repo_user.Update{}),
			repo_user.Filter{Username: &toUsername}).Return(nil)

		u.transferHistoryName.EXPECT().Create(gomock.Any(), models.TransferHistoryName{
			SenderName:   fromUsername,
			ReceiverName: toUsername,
			Amount:       transferAmount,
		}).Return(uuid.New(), nil)

		u.mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

		err := u.self.TransferCoins(ctx, toUsername, transferAmount)
		u.NoError(err)
	})

	u.Run("недостаточно средств", func() {
		insufficientAmount := uint(2000)

		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(fromUserID, nil)
		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{Username: &toUsername}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Name: toUsername, Balance: initialBalance}, nil)
		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{ID: &fromUserID}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Name: fromUsername, Balance: initialBalance}, nil)

		err := u.self.TransferCoins(ctx, toUsername, insufficientAmount)
		u.Error(err)
		u.Contains(err.Error(), "no enough money")
	})
}

func (u *UsecaseSuite) TestUsecase_Buy() {
	ctx := context.Background()
	userID := uuid.New()
	itemName := "test_item"
	initialBalance := uint(1000)
	itemPrice := uint(500)
	itemID := uuid.New()

	u.Run("успешная покупка предмета", func() {
		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(userID, nil)

		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{ID: &userID}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Balance: initialBalance}, nil)
		u.itemRepo.EXPECT().Get(gomock.Any(), repo_item.Filter{Name: &itemName}).
			Return(models.Item{ID: itemID, Price: itemPrice}, nil)

		u.userRepo.EXPECT().Update(gomock.Any(),
			gomock.AssignableToTypeOf(repo_user.Update{}),
			repo_user.Filter{ID: &userID}).Return(nil)

		u.userItemRepo.EXPECT().Create(gomock.Any(), models.UserItem{
			UserID: userID,
			ItemID: itemID,
		}).Return(uuid.New(), nil)

		u.mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

		err := u.self.Buy(ctx, itemName)
		u.NoError(err)
	})

	u.Run("недостаточно средств для покупки", func() {
		expensiveItemPrice := uint(2000)

		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(userID, nil)
		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{ID: &userID}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Balance: initialBalance}, nil)
		u.itemRepo.EXPECT().Get(gomock.Any(), repo_item.Filter{Name: &itemName}).
			Return(models.Item{Price: expensiveItemPrice}, nil)

		err := u.self.Buy(ctx, itemName)
		u.Error(err)
		u.Contains(err.Error(), "no enough money")
	})

	u.Run("предмет не найден", func() {
		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(userID, nil)
		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.userRepo.EXPECT().Get(gomock.Any(), repo_user.Filter{ID: &userID}, repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{Balance: initialBalance}, nil)
		u.itemRepo.EXPECT().Get(gomock.Any(), repo_item.Filter{Name: &itemName}).
			Return(models.Item{}, common.ErrNotFound)

		err := u.self.Buy(ctx, itemName)
		u.Error(err)
		u.Contains(err.Error(), "item")
	})
}

func (u *UsecaseSuite) TestUsecase_Info() {
	ctx := context.Background()
	userID := uuid.New()
	userName := "test_user"
	balance := uint(1000)

	u.Run("успешное получение информации", func() {
		u.userRepo.EXPECT().IsAuth(gomock.Any()).Return(userID, nil)

		u.transactionFabric.EXPECT().Begin(gomock.Any()).Return(ctx, u.mockTx, nil)

		u.mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

		u.userRepo.EXPECT().Get(gomock.Any(),
			repo_user.Filter{ID: &userID},
			repo_user.GetOptions{ForUpdate: true}).
			Return(models.User{ID: userID, Name: userName, Balance: balance}, nil)

		itemID := uuid.New()
		userItemsAmount := []models.UserItemsAmount{{
			ItemID:   itemID,
			Quantity: 1,
		}}
		u.userItemRepo.EXPECT().GetUserItemsAmount(gomock.Any(),
			repo_useritem.Filter{UserID: &userID}).
			Return(userItemsAmount, nil)

		items := []models.Item{{
			ID:    itemID,
			Name:  "test_item",
			Price: 100,
		}}
		u.itemRepo.EXPECT().GetMany(gomock.Any(), repo_item.Filter{}).
			Return(items, nil)

		receivedHistory := []models.TransferHistoryName{{
			SenderName:   "sender",
			ReceiverName: userName,
			Amount:       200,
		}}
		u.transferHistoryName.EXPECT().GetMany(gomock.Any(),
			repo_transferhistory.Filter{ReceiverName: &userName}).
			Return(receivedHistory, nil)

		sentHistory := []models.TransferHistoryName{{
			SenderName:   userName,
			ReceiverName: "receiver",
			Amount:       100,
		}}
		u.transferHistoryName.EXPECT().GetMany(gomock.Any(),
			repo_transferhistory.Filter{SenderName: &userName}).
			Return(sentHistory, nil)

		u.mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

		info, err := u.self.Info(ctx)
		u.NoError(err)
		u.Equal(balance, info.Coins)
		u.Len(info.Inventory, 1)
		u.Len(info.CoinHistory.Received, 1)
		u.Len(info.CoinHistory.Sent, 1)
	})

	u.Run("ошибка авторизации", func() {
		u.userRepo.EXPECT().IsAuth(gomock.Any()).
			Return(uuid.UUID{}, common.ErrUnauthorized)

		info, err := u.self.Info(ctx)
		u.Error(err)
		u.Equal(err, common.ErrUnauthorized)
		u.Equal(models.UserInfo{}, info)
	})
}
