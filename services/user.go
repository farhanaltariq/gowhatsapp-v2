package services

import (
	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	InsertOrUpdate(models.User) (models.User, error)
	GetDataByUsernameOrEmail(user models.User) (models.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db}
}

func (server *userService) InsertOrUpdate(user models.User) (models.User, error) {
	query := server.db.Debug()
	err := query.Save(&user).Error
	if err != nil {
		logrus.Errorln(err)
		return user, err
	}

	return user, nil
}

func (server *userService) GetDataByUsernameOrEmail(user models.User) (models.User, error) {
	var data models.User
	err := server.db.Debug().Where("username = ? OR email = ?", user.Username, user.Email).Where("deleted_at IS NULL").Where("id != ? ", user.ID).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return data, nil
		}
		logrus.Errorln(err)
		return data, err
	}

	return data, nil
}
