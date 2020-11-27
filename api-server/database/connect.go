package database

import (
	"fmt"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/models"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

// Connect with Postgres DB ...
func ConnectDB() {
	dbPort := config.GetConfig("DB_PORT")
	port, err := strconv.ParseUint(dbPort, 10, 32)
	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.GetConfig("DB_HOST"), port, config.GetConfig("DB_USER"), config.GetConfig("DB_PASSWORD"), config.GetConfig("DB_NAME")))
	if err != nil {
		log.Fatal("failed to connect database")
	}
	DB.AutoMigrate(&models.User{}, &models.Client{}, &models.Field{}, &models.Template{}, &models.Job{}, &models.Auth{})
	DB.Model(&models.Client{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	DB.Model(&models.Field{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")
	DB.Model(&models.Template{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	DB.Model(&models.Job{}).AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")
	DB.Model(&models.Auth{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	repository.DB = DB
	log.Println("Database connection successful")
}
