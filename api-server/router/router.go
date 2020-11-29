package router

import (
	"github.com/Harry-027/go-notify/api-server/handler"
	"github.com/Harry-027/go-notify/api-server/middleware"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Auth
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/signup", handler.Signup)
	auth.Post("/forgotPassword", handler.ForgotPassword)
	auth.Post("/getNewPassword/:id", handler.GetNewPassword)

	//Privacy
	privacy := app.Group("/privacy", middleware.Protected())
	privacy.Post("/updatePassword", handler.UpdatePassword)
	privacy.Delete("/deleteAccount", handler.DeleteAccount)
	privacy.Post("/logout", handler.Logout)

	// Users Apis
	user := app.Group("/api", middleware.Protected())
	user.Get("/users", handler.GetUsers)
	user.Get("/subscriptionDetails", handler.GetSubsDetail)
	user.Get("/clients", handler.GetClients)
	user.Get("/clientDetails/:clientID", handler.GetClientVariable)
	user.Post("/clients", handler.AddClients)
	user.Post("/subscribe", handler.Subscribe)
	user.Put("/updateClient/:mailId", handler.UpdateClient)
	user.Put("/updateTemplate/:id", handler.UpdateTemplate)
	user.Post("/addTemplate", handler.AddTemplate)
	user.Get("/getTemplates", handler.GetTemplates)
	user.Post("/clientDetails", handler.AddUserVariable)
	user.Delete("/clientDetails", handler.DeleteUserVariable)
	user.Delete("/deleteClient/:mailId", handler.DeleteClient)
	user.Delete("/deleteTemplate/:id", handler.DeleteTemplate)

	// sendMail
	user.Post("/sendMail", handler.SendMail)
	user.Post("/scheduleMail", handler.ScheduleMail)
	user.Post("/deleteScheduleMail", handler.DeleteScheduleMail)
	user.Get("/checkAuditLog", handler.CheckAuditLog)
	// swagger
	app.Get("/swagger/*", swagger.Handler) // default
}
