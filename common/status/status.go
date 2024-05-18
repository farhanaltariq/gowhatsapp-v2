package status

import (
	"fmt"

	"github.com/farhanaltariq/fiberplate/common"
	"github.com/gofiber/fiber/v2"
)

func Errorf(c *fiber.Ctx, codes int, message string, args ...interface{}) error {
	if codes >= 200 && codes < 300 {
		return Successf(c, codes, message, args...)
	}

	jsonMsg := &common.ResponseMessage{
		Error:   true,
		Code:    codes,
		Message: fmt.Sprintf(message, args...),
	}
	return c.Status(codes).JSON(jsonMsg)
}

func Successf(c *fiber.Ctx, codes int, message string, args ...interface{}) error {
	if !(codes >= 200 && codes < 300) {
		return Errorf(c, codes, message, args...)
	}
	jsonMsg := &common.ResponseMessage{
		Error:   false,
		Code:    codes,
		Message: fmt.Sprintf(message, args...),
	}

	return c.Status(codes).JSON(jsonMsg)
}
