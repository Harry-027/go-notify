package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect with Postgres DB ...
func ConnectDB() {
	dbPort := config.GetConfig(config.DB_PORT)
	port, err := strconv.ParseUint(dbPort, 10, 32)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.GetConfig("DB_HOST"), port, config.GetConfig("DB_USER"), config.GetConfig("DB_PASSWORD"), config.GetConfig("DB_NAME"))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Audit{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Auth{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Client{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Job{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Field{})
	ErrMigrate(err)

	err = db.AutoMigrate(&models.Template{})
	ErrMigrate(err)

	repository.DB = db
	log.Println("Database connection successful")
}

func ErrMigrate(err error) {
	if err != nil {
		log.Fatal("failed to migrate database")
	}
}
