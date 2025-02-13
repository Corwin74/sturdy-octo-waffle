// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"shop/internal/conf"
	"shop/internal/repository/user"
	"shop/internal/server"
	"shop/internal/service"
	user2 "shop/internal/usecase/user"
	"shop/pkg/querier"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, data *conf.Data, secrets *conf.Secrets, logger log.Logger) (*kratos.App, func(), error) {
	database, err := querier.NewDatabase(data)
	if err != nil {
		return nil, nil, err
	}
	repository := user.NewRepository(database)
	usecase := user2.NewUsecase(repository, secrets)
	shopService := service.NewShopService(usecase)
	grpcServer := server.NewGRPCServer(confServer, shopService, logger)
	httpServer := server.NewHTTPServer(confServer, shopService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
	}, nil
}
