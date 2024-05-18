package routes

import (
	"github.com/farhanaltariq/fiberplate/controllers"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/gofiber/fiber/v2"
)

func Authentications(router fiber.Router, service middleware.Services) {
	authController := controllers.NewAuthController(service)

	router.Post("/login", authController.Login)
	router.Post("/register", authController.Register)
}
