package user

import (
	"shop/internal/conf"
	"shop/mocks/pkg/querier"
	"shop/mocks/pkg/transaction"
	"shop/mocks/usecase"
	"testing"

	"github.com/golang/mock/gomock"
) 

func TestUsecase_SomeMethod(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Создаем моки
    mockUserRepo := usecase.NewMockUserRepo(ctrl)
    mockItemRepo := usecase.NewMockItemRepo(ctrl)
    mockUserItemRepo := usecase.NewMockUserItemRepo(ctrl)
    mockTransferHistory := usecase.NewMockTransferHistory(ctrl)
    mockTransferHistoryName := usecase.NewMockTransferHistoryName(ctrl)
    mockQuerier := querier.NewMockQuerier(ctrl)
    mockTransactionFabric := transaction.NewMockFabric(ctrl)

    // Создаем тестовый конфиг
    testConfig := &conf.Secrets{
		JwtKey: "some-secret-ley",
        // заполняем нужными тестовыми значениями
    }

    // Создаем экземпляр usecase с моками
    usecase := &Usecase{
        userRepo:             mockUserRepo,
        itemRepo:             mockItemRepo,
        userItemRepo:         mockUserItemRepo,
        transferHistory:      mockTransferHistory,
        transferHistoryName:  mockTransferHistoryName,
        config:              testConfig,
        querier:             mockQuerier,
        transactionFabric:   mockTransactionFabric,
    }

    // // Настраиваем ожидаемое поведение моков
    mockUserRepo.EXPECT().Get().Return(expectedResult, nil)

    // // Выполняем тестируемый метод
    // result, err := usecase.SomeMethod()

    // // Проверяем результаты
    // assert.NoError(t, err)
    // assert.Equal(t, expectedResult, result)
	usecase.userItemRepo.Get()
}