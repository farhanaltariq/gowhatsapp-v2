package aiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/farhanaltariq/fiberplate/database/models"
	"github.com/farhanaltariq/fiberplate/utils"
)

func AskChatGPT(q *string) (*models.GPTResponse, error) {
	url := utils.GetEnv("OPENAI_URL", "https://api.openai.com/v1/chat/completions")
	apiKey := utils.GetEnv("OPENAI_API_KEY", "secret")
	res := &models.GPTResponse{}

	payload := models.GPTRequest{
		Model:       utils.GetEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		Messages:    []models.GPTMessage{{Role: "user", Content: *q}},
		Temperature: 0.7,
	}
	postBody, _ := json.Marshal(payload)
	jsonPayload := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", url, jsonPayload)
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, nil
}
