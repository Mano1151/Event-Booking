//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"notificationservice/internal/biz"
	"notificationservice/internal/conf"
	"notificationservice/internal/data"
	"notificationservice/internal/server"
	"notificationservice/internal/service"
)

func wireApp(confServer *conf.Server, confData *conf.Data, confEmail *conf.Email, logger log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		// data providers
		data.NewDB,        // -> *gorm.DB, cleanup
		data.NewData,      // -> *Data, cleanup
		data.NewNotificationRepo,

		// clients
		data.ProvideBookingClient,
		data.ProvideUserClient,

		// biz providers
		biz.ProvideEmailConfig,
		biz.NewNotificationUsecase,

		// service + server + app
		service.NewNotificationService,
		server.NewGRPCServer,
		server.NewHTTPServer,
		newApp,
	)
	return nil, nil, nil
}
