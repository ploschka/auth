package model

import (
	"database/sql"
	"os"

	log "github.com/ploschka/auth/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Guid       string `gorm:"primaryKey"`
	Admin      bool
	Email      string
	RefreshKey sql.NullString
}

var db *gorm.DB

func init() {
	log.Info("Model init start")
	link, ok := os.LookupEnv("DATABASE_LINK")
	if !ok {
		panic("DATABASE_LINK environment variable is undefined")
	}
	var err error
	db, err = gorm.Open(postgres.Open(link), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
	log.Info("Model init ended")
}

func GetDB() *gorm.DB {
	return db
}
