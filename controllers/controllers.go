package controllers

import (
	"github.com/farhanaltariq/fiberplate/common/codes"
	"github.com/farhanaltariq/fiberplate/common/status"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/gofiber/fiber/v2"
)

type MiscController interface {
	HealthCheck(c *fiber.Ctx) error
}

type controller struct {
	middleware.Services
}

func NewMiscController(service middleware.Services) MiscController {
	return &controller{service}
}
func (server *controller) HealthCheck(c *fiber.Ctx) error {
	return status.Successf(c, codes.OK, "Server Running")
}
