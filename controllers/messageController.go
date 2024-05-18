package controllers

import (
	"encoding/json"

	"github.com/farhanaltariq/fiberplate/common/codes"
	"github.com/farhanaltariq/fiberplate/common/status"
	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/farhanaltariq/fiberplate/libs/whatsapp"
	"github.com/farhanaltariq/fiberplate/middleware"
	"github.com/gofiber/fiber/v2"
)

type MessageController interface {
	SendMessage(c *fiber.Ctx) error
}

func NewMessageController(service middleware.Services) MessageController {
	return &controller{service}
}

// @Summary Send Message
// @Description Send New Message to Desired Number
// @Tags Message
// @Accept json
// @Param data body models.Message true "Message data"
// @Produce json
// @Success 200 {object} common.ResponseMessage
// @Failure 400 {object} common.ResponseMessage
// @Router /message [post]
func (s *controller) SendMessage(c *fiber.Ctx) error {
	messageData := &models.Message{}

	if err := json.Unmarshal(c.Body(), &messageData); err != nil {
		return status.Errorf(c, codes.BadRequest, err.Error())
	}

	err := whatsapp.SendMessage(s.WaClient, messageData.Number, messageData.Message)
	if err != nil {
		return status.Errorf(c, codes.BadRequest, "Failed to send message")
	}
	return status.Successf(c, codes.OK, "Success")
}
