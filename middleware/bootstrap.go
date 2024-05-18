package middleware

import (
	"github.com/farhanaltariq/fiberplate/database"
	"github.com/farhanaltariq/fiberplate/libs/whatsapp"
	"github.com/farhanaltariq/fiberplate/services"
	"go.mau.fi/whatsmeow"
	"gorm.io/gorm"
)

type Services struct {
	DB          *gorm.DB
	AuthService services.AuthenticationService
	UserService services.UserService
	WaClient    *whatsmeow.Client
}

func InitServices() Services {
	db := database.GetDBConnection()
	waClient := whatsapp.GetClient()
	return Services{
		DB:          db,
		AuthService: services.NewAuthService(db),
		UserService: services.NewUserService(db),
		WaClient:    waClient,
	}
}
