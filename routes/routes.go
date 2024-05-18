package routes

import (
	"github.com/farhanaltariq/fiberplate/controllers"
	"github.com/farhanaltariq/fiberplate/database"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/farhanaltariq/fiberplate/services"
	"github.com/gofiber/fiber/v2"
)

func initServices() middleware.Services {
	db := database.GetDBConnection()
	return middleware.Services{
		DB:          db,
		AuthService: services.NewAuthService(db),
		UserService: services.NewUserService(db),
	}
}

func Init(app *fiber.App) {
	services := initServices()

	api := app.Group("/api")
	api.Get("/", controllers.NewMiscController(services).HealthCheck)

	Authentications(api.Group("/auth"), services)

	api.Use(middleware.AuthInterceptor)
}
