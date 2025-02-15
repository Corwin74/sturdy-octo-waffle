package service

import (
	"shop/internal/repository/transferhistory"
	"shop/internal/repository/user"
	userUsecase "shop/internal/usecase/user"
	"shop/pkg/querier"
	"shop/pkg/transaction"

	"github.com/google/wire"
)

// ProviderSet is service providers.

// ShopServiceSet -- собираем все вместе
var ShopServiceSet = wire.NewSet(
	querier.NewDatabase,
	transaction.NewTrFabric,
	user.NewRepository,
	transferhistory.NewRepository,
	userUsecase.NewUsecase,
	
	NewShopService,
	// Привязываем интерфейсы к реализации
	//wire.Bind(new(transaction.Fabric), new(*transaction.TrFabricImpl)),
	wire.Bind(new(querier.Querier), new(*querier.Database)),
	// repos
	
	wire.Bind(new(userUsecase.UserRepo), new(*user.Repository)),
	wire.Bind(new(userUsecase.TransferHistory), new(*transferhistory.Repository)),
	// usecase
	wire.Bind(new(UserUsecase), new(*userUsecase.Usecase)),
	// wire.Bind(new(Tr), new(*userUsecase.Usecase)),
)
