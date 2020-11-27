package middleware

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"os"
	"time"
)

// Set up the required middlewares before server startup ...
func SetUpMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		Max: 50,
	}))
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Next:         nil,
		Format:       "[${time}] ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stderr,
	}))
	prom := prometheusMiddleware(app)
	app.Use(prom.Middleware)
}

func prometheusMiddleware(app *fiber.App) *fiberprometheus.FiberPrometheus {
	prometheus := fiberprometheus.New("prometheus-service")
	prometheus.RegisterAt(app, "/metrics")
	return prometheus
}
