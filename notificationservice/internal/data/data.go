package data

import (
	
	"notificationservice/internal/conf"
     "github.com/google/wire"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)




// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(ProvideBookingClient,NewData,NewDB,NewNotificationRepo,ProvideUserClient)

// Data holds connections (DB etc.)
type Data struct {
	db *gorm.DB
}

// NewDB initializes *gorm.DB and returns a cleanup func.
func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, func(), error) {
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// Auto-migrate Notification table
	if err := db.AutoMigrate(&Notification{}); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		log.NewHelper(logger).Info("closing the database resources")
	}

	return db, cleanup, nil
}

// NewData returns a Data wrapper object (so other providers can accept *Data)
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	d := &Data{db: db}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return d, cleanup, nil
}

// ---------------- NotificationRepo ----------------

// Notification is the DB model for notifications
