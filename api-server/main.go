package main

import (
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/database"
	_ "github.com/Harry-027/go-notify/api-server/docs"
	"github.com/Harry-027/go-notify/api-server/handler"
	"github.com/Harry-027/go-notify/api-server/middleware"
	"github.com/Harry-027/go-notify/api-server/router"
	"github.com/Harry-027/go-notify/api-server/utils"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func main() {
	app := fiber.New() // new fiber app
	pipelineSetup(app) // configure the service pipeline
	portNumber := config.GetConfig(config.SERVER_PORT)
	err := app.Listen(portNumber) // start the app
	if err != nil {
		log.Fatal("App crashed :: An error occurred", err.Error())
	}
}

// pipeline setup ...
func pipelineSetup(app *fiber.App) {
	config.LoadConfig()              // load env variables
	database.ConnectDB()             // connect to postgres db
	middleware.SetUpMiddlewares(app) // setup middlewares
	router.SetupRoutes(app)          // setup routes
	conn := utils.InitKafkaConn()    // connect to kafka
	handler.KafkaConn = conn
	handler.RedisPoolInit()
}
