package middleware

import (
	"github.com/farhanaltariq/fiberplate/services"
	"gorm.io/gorm"
)

type Services struct {
	DB          *gorm.DB
	AuthService services.AuthenticationService
	UserService services.UserService
}
