package routes

import (
	"github.com/farhanaltariq/fiberplate/controllers"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/gofiber/fiber/v2"
)

func Messages(router fiber.Router, service middleware.Services) {
	messageController := controllers.NewMessageController(service)

	router.Post("/", messageController.SendMessage)
	router.Post("/logout", messageController.Logout)
	router.Get("/qr", messageController.GenerateQR)
}
