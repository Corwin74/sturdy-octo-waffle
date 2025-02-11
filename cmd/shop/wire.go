//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"shop/internal/conf"
	"shop/internal/server"
	"shop/internal/service"
	"shop/internal/usecase"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(usecase.ProviderSet, server.ProviderSet, service.ProviderSet, newApp))
}
