package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	UserType string `json:"userType"`
}
