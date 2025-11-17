// +build wireinject

package main

import (
	"paymentservice/internal/biz"
	"paymentservice/internal/conf"
	"paymentservice/internal/data"
	"paymentservice/internal/server"
	"paymentservice/internal/service"

	"github.com/google/wire"
	"github.com/go-kratos/kratos/v2"       // <<--- Import kratos
	"github.com/go-kratos/kratos/v2/log"
)

// wireApp initializes PaymentService app.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		data.ProviderSet,            // Provides DB, PaymentRepo, BookingClient
		biz.NewPaymentUsecase,       // Injects PaymentRepo + BookingClient
		service.NewPaymentService,   // Injects PaymentUsecase
		server.NewGRPCServer,        // Injects PaymentService + confServer + logger
		server.NewHTTPServer,        // Injects PaymentService + confServer + logger
		newApp,                      // Combines GRPC + HTTP servers
	)
	return nil, nil, nil
}
