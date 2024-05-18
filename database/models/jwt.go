package models

import "github.com/golang-jwt/jwt/v4"

type JWT struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   int    `json:"userId"`
}

type Claims struct {
	Data JWT `json:"data"`
	jwt.RegisteredClaims
}
