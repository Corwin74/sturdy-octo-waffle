package service

import (
	"shop/pkg/querier"
	"shop/internal/repository/user"
	userUsecase "shop/internal/usecase/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ShopServiceSet = wire.NewSet(
	querier.NewDatabase,
	user.NewRepository,
	userUsecase.NewUsecase,
	NewShopService,
	// Привязываем интерфейсы к реализации
	wire.Bind(new(querier.Querier), new(*querier.Database)),
	wire.Bind(new(userUsecase.UserRepo), new(*user.Repository)),
	wire.Bind(new(UserUsecase), new(*userUsecase.Usecase)),
)
