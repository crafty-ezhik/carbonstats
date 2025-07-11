package db

import (
	"fmt"
	"github.com/crafty-ezhik/carbonstats/config"
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func GetConnection(config *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(25)                 // Количество "холостых" соединений.
	sqlDB.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Время жизни соединения
	return db
}

func GoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		statistics.ClientStatistics{},
		service_description.ServiceDescription{},
	)
	if err != nil {
		panic(err)
	}
}
