package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/farhanaltariq/fiberplate/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseGPT struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type savedMessage struct {
	data map[interface{}]string
}

var SavedMessage = savedMessage{
	data: make(map[interface{}]string),
}

func hitAI(msg string, sender interface{}, responseChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	if msg == "/reset" {
		SavedMessage.data[sender] = ""
		responseChan <- "Successfully reset the conversation, you can start a new conversation now"
		return
	}

	url := "https://api.openai.com/v1/chat/completions"
	apiKey := utils.GetEnv("OPENAI_API_KEY", "secret")

	reqMsg := SavedMessage.data[sender] + "]\n" + msg
	if SavedMessage.data[sender] != "" {
		SavedMessage.data[sender] = SavedMessage.data[sender] + "\nmessage : " + msg
	} else {
		SavedMessage.data[sender] = "\nprevious_messages : [\nmessage : " + msg
		reqMsg = msg
	}

	fmt.Print("\033[31m\n", reqMsg, "]\033[0m\n\n")

	payload := Request{
		Model:       "gpt-3.5-turbo",
		Messages:    []Message{{Role: "user", Content: reqMsg}},
		Temperature: 0.7,
	}
	postBody, _ := json.Marshal(payload)
	jsonPayload := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", url, jsonPayload)
	if err != nil {
		fmt.Println(err)
		responseChan <- "" // Send empty response on error
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		responseChan <- "" // Send empty response on error
		return
	}
	defer resp.Body.Close()

	response := ResponseGPT{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(err)
		responseChan <- "" // Send empty response on error
		return
	}

	if len(response.Choices) < 1 {
		// say if can't find any response
		responseChan <- "Sorry, I don't understand what you mean"
		return
	}

	msg = response.Choices[0].Message.Content
	fmt.Println("\033[34m", msg, "\033[0m")

	SavedMessage.data[sender] = SavedMessage.data[sender] + "\nanswer : " + msg + "\n"

	responseChan <- msg // Send the AI response
}

func EventHandler(client *whatsmeow.Client, evt interface{}, debug bool) {
	switch v := evt.(type) {
	case *events.Message:
		// if not from me and not empty
		msg := v.Message.ExtendedTextMessage.GetText()
		if msg == "" {
			msg = v.Message.GetConversation()
		}

		sender := v.Info.MessageSource.Chat
		senderName := v.Info.PushName
		if v.Info.IsFromMe || msg == "" {
			return
		}

		fmt.Println("\033[32mSender\t:", senderName, " | ", sender, "\033[0m")
		fmt.Println("\033[32mMessage\t:", msg, "\033[0m")
		fmt.Println("\033[32mReply\t:", SavedMessage.data[sender], "\033[0m")

		if debug {
			return
		}

		responseChan := make(chan string, 1)
		var wg sync.WaitGroup
		wg.Add(1)

		go hitAI(msg, sender, responseChan, &wg)

		// Wait for the AI response concurrently
		go func() {
			wg.Wait() // Wait for the AI processing to complete
			close(responseChan)
			msg, ok := <-responseChan
			if ok && msg != "" {
				protoMsg := &proto.Message{
					ExtendedTextMessage: &proto.ExtendedTextMessage{
						// text to be sent to the sender
						Text: &msg,
					},
				}
				client.SendMessage(context.Background(), sender, protoMsg)
			}
		}()
	}
}
