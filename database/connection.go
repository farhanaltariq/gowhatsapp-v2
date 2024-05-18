package database

import (
	"fmt"

	model "github.com/farhanaltariq/fiberplate/database/models"
	"github.com/farhanaltariq/fiberplate/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Database struct {
	*gorm.DB
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Authentications{},
	)
}

func Connect() error {
	dsn := utils.GetEnv("DB_DSN", "user=user password=password dbname=fiber port=5432 sslmode=disable TimeZone=Asia/Jakarta")
	var config gorm.Dialector

	switch utils.GetEnv("DB_DRIVER", "postgres") {
	case "postgres":
		config = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: false, // disables implicit prepared statement usage
		})
	case "mysql":
		config = mysql.New(mysql.Config{
			DSN: dsn,
		})
	default:
		logrus.Errorln("Error connecting to database : driver not set")
		return fmt.Errorf("driver not set")
	}

	db, err := gorm.Open(config, &gorm.Config{})
	if err != nil {
		logrus.Errorln("Error connecting to database", err)
		return err
	}

	AutoMigrate(db)
	SetupDBConnection(db)
	return nil
}

func SetupDBConnection(DB *gorm.DB) {
	db = DB
}

func GetDBConnection() *gorm.DB {
	return db
}
