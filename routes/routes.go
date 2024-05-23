package routes

import (
	"github.com/farhanaltariq/fiberplate/controllers"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Init(app *fiber.App) {
	services := middleware.InitServices()

	app.Use(cors.New())

	api := app.Group("/api")
	api.Get("/", controllers.NewMiscController(services).HealthCheck)

	// Reminders : Move below authentication route
	Messages(api.Group("/message"), services)
	// End Reminders

	Authentications(api.Group("/auth"), services)

	api.Use(middleware.AuthInterceptor)
}
