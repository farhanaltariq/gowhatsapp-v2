package utils

import (
	"encoding/base64"
	"fmt"
)

func Decrypt(password string, salt string) (string, error) {
	// signature := GetEnv("SIGNATURE_KEY", "secret")

	plainEnc, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", fmt.Errorf("%s", "Failed to decode password")
	}

	sigIndex := 0
	pass := []rune(string(plainEnc))
	for i := range pass {
		if sigIndex == len(salt) {
			sigIndex = 0
		}

		pass[i] -= rune(salt[sigIndex])

		sigIndex++
	}

	sigIndex = 0
	decodedPass := make([]rune, 0)
	for i := range pass {
		if i%2 == 0 {
			decodedPass = append(decodedPass, pass[i])
		}
	}

	return string(decodedPass), nil
}
