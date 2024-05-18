package models

import "gorm.io/gorm"

type Authentications struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primaryKey"`
	UserId   int    `json:"userId"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Register struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Country         string `json:"country"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type Login struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type AuthenticationResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiredAt    string `json:"expired_at"`
}
