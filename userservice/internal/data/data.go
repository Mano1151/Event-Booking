package data

import (
	"fmt"
	"userservice/internal/conf"
	
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewUserRepo)

// Data holds DB client
type Data struct {
	db *gorm.DB
}

// NewData initializes Data with GORM DB.
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	// ✅ Run AutoMigrate
	err := db.AutoMigrate(
		&User{},
		
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	helper.Info("✅ AutoMigrate completed")

	cleanup := func() {
		helper.Info("closing the database connection")
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	return &Data{db: db}, cleanup, nil
}

// NewDB opens a Postgres connection with GORM
func NewDB(c *conf.Data) (*gorm.DB, func(), error) {
    dsn := c.Database.Source
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, nil, err
    }

    // Connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    cleanup := func() {
        _ = sqlDB.Close()
    }

    return db, cleanup, nil
}
