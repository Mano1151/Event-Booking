package data

import (
	"context"
	
	"paymentservice/internal/biz"
	"paymentservice/internal/conf"

	bookingv1 "bookingservice/api/bookingservice/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet wires data dependencies
var ProviderSet = wire.NewSet(
	ProvideDSN,
	NewDB,
	ProvideBookingClient,
	NewPaymentRepo,
	NewBookingClient,
)

type Data struct {
	PaymentRepo biz.PaymentRepo
}

// ProvideDSN returns the database DSN string from config
func ProvideDSN(confData *conf.Data) string {
    return confData.Database.Source
}

// NewBookingClient wraps gRPC client into your interface
func NewBookingClient(c bookingv1.BookingServiceClient) biz.BookingClient {
	return &bookingClient{client: c} // bookingClient implements biz.BookingClient
}

// NewData initializes Data struct
func NewData(paymentRepo biz.PaymentRepo) (*Data, func(), error) {
	cleanup := func() {}
	return &Data{PaymentRepo: paymentRepo}, cleanup, nil
}

// NewDB creates a GORM DB connection
func NewDB(dsn string) (*gorm.DB, func(), error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 👉 Run AutoMigrate here
	if err := db.AutoMigrate(&PaymentModel{}); err != nil {
		return nil, nil, err
	}

	cleanup := func() { _ = sqlDB.Close() }
	return db, cleanup, nil
}


// ProvideBookingClient creates a gRPC BookingService client
func ProvideBookingClient() (bookingv1.BookingServiceClient, func(), error) {
	ctx := context.Background()
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint("127.0.0.1:9002"))
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = conn.Close() }
	client := bookingv1.NewBookingServiceClient(conn)
	return client, cleanup, nil
}
