package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// load env variables ...
func LoadConfig() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// get the given env variable value
func GetConfig(key string) string {
	return os.Getenv(key)
}
