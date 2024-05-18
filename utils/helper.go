package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Method = map[string]string{
	"DELETE":  "\033[31m",
	"HEAD":    "\033[32m",
	"GET":     "\033[32m",
	"POST":    "\033[33m",
	"PUT":     "\033[34m",
	"PATCH":   "\033[35m",
	"OPTIONS": "\033[36m",
	"RESET":   "\033[0m",
}

func FormatMethod(c *fiber.Ctx) string {
	method := c.Method()
	if method == "HEAD" || method == "GET" || method == "POST" || method == "PUT" || method == "PATCH" {
		return Method[method] + method + Method["RESET"] + "\t\t"
	}
	return Method[method] + method + Method["RESET"] + "\t"
}

func GetEnv(key, fallback string) string {

	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			logrus.Errorln("Failed to load .env", err)
			return fallback
		}
	}

	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func CustomFormatter() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05"
	logrus.SetFormatter(customFormatter)

	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
}

func GetExpirationTime() time.Time {
	expiryTime := GetEnv("JWT_EXPIRY_TIME", "30m")
	// get latest index
	formatTime := expiryTime[len(expiryTime)-1:]
	// handle 1m, 1h, 1d

	timeEnv := expiryTime[:len(expiryTime)-1]
	duration, err := strconv.Atoi(timeEnv)
	if err != nil {
		logrus.Errorln(err)
		return time.Now().Add(time.Minute * 30)
	}

	switch formatTime {
	case "m":
		return time.Now().Add(time.Minute * time.Duration(duration))
	case "h":
		return time.Now().Add(time.Hour * time.Duration(duration))
	case "d":
		return time.Now().Add(time.Hour * 24 * time.Duration(duration))
	}

	logrus.Infoln(formatTime)
	return time.Now().Add(time.Minute * 3)
}

func GenerateToken(user *models.User) (string, error) {
	GetExpirationTime()
	key := []byte(GetEnv("JWT_SECRET", "secret"))
	accessToken := models.JWT{
		Username: user.Username,
		Email:    user.Email,
		UserId:   user.ID,
	}

	expiredAt := time.Now().Add(time.Minute * 3).Unix()
	numericDate := jwt.NewNumericDate(time.Unix(expiredAt, 0))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		Data: accessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: numericDate,
		},
	})
	signed, err := token.SignedString([]byte(key))
	if err != nil {
		logrus.Errorln(err)
		return "", err
	}
	return signed, nil
}

func RenameBaseUrlInFile(relativePath, baseUrl string) {
	// Get the absolute path
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	// Open the file
	file, err := os.Open(absolutePath)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	defer file.Close()

	// Read the file
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	// Replace base URL
	newData := []byte(strings.ReplaceAll(string(data), "localhost:3000", baseUrl))

	// Write the updated data back to the file
	err = os.WriteFile(absolutePath, newData, 0644)
	if err != nil {
		logrus.Errorln(err)
		return
	}
}

func RenameBaseUrlSwagger(baseUrl string) {
	// Rename base URL in swagger.json
	RenameBaseUrlInFile("docs/swagger.json", baseUrl)

	// Rename base URL in docs.go
	RenameBaseUrlInFile("docs/docs.go", baseUrl)

	// Rename base URL in swagger.yaml
	RenameBaseUrlInFile("docs/swagger.yaml", baseUrl)
}
