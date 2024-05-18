package services

import (
	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	InsertOrUpdate(models.Authentications) error
	GetDataByUserId(int) (models.Authentications, error)
}

type authenticationService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthenticationService {
	return &authenticationService{db}
}
func (server *authenticationService) InsertOrUpdate(auth models.Authentications) error {
	query := server.db.Debug()
	err := query.Save(&auth).Error
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

func (server *authenticationService) GetDataByUserId(id int) (models.Authentications, error) {
	var data models.Authentications
	err := server.db.Debug().Where("user_id = ?", id).Where("deleted_at IS NULL").First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		logrus.Errorln(err)
		return data, err
	}

	return data, nil
}
