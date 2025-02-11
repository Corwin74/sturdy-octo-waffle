package service

import (
	"shop/internal/usecase/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(wire.Bind(new(UserUsecase), new(*user.Usecase)), NewShopService)
