//go:build wireinject
// +build wireinject

package main

import (
	"bookingservice/internal/biz"
	"bookingservice/internal/conf"
	"bookingservice/internal/data"
	"bookingservice/internal/server"
	"bookingservice/internal/service"

	

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
