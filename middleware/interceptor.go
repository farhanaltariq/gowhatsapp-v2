package middleware

import (
	"fmt"
	"strings"

	"github.com/farhanaltariq/fiberplate/common/codes"
	"github.com/farhanaltariq/fiberplate/common/status"
	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/farhanaltariq/fiberplate/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func CommonMiddleware(c *fiber.Ctx) error {
	endpoint := c.BaseURL() + c.Path()
	logrus.Infoln(utils.FormatMethod(c), endpoint)
	return c.Next()
}

func validateToken(tokenString string, jwtSecret []byte) error {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return err
	}

	_, ok := token.Claims.(*models.Claims)
	if !ok {
		//lint:ignore ST1005 will sent to user
		return fmt.Errorf("Invalid token")
	}

	return nil
}

func AuthInterceptor(c *fiber.Ctx) error {
	jwtSecret := []byte(utils.GetEnv("JWT_SECRET", "secret"))
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return status.Errorf(c, codes.Unauthorized, "Unauthorized")
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix

	if err := validateToken(tokenString, jwtSecret); err != nil {
		return status.Errorf(c, codes.Unauthorized, "Unauthorized")
	}

	return c.Next()
}
