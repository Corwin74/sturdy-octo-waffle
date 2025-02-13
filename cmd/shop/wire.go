//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"shop/internal/conf"
	"shop/internal/server"
	"shop/internal/service"
	//"shop/internal/repository/user"
	//userUsecase "shop/internal/usecase/user"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(serverConf *conf.Server, data *conf.Data, secrets *conf.Secrets, logger log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		wire.ProviderSet(wire.Bind(new(*conf.Secrets), secrets)),
		// Data providers
		// ...
		// Repository providers
		//user.NewRepository,
		// Usecase providers
		//userUsecase.NewUsecase,
		
		//Service providers
		//service.NewShopService,
		//Server providers
		server.NewGRPCServer,
        server.NewHTTPServer,
		// App providers
		newApp,
		// Wire providers set
		service.ShopServiceSet,
	))

}
