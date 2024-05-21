package whatsapp

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/farhanaltariq/fiberplate/libs/aiclient"
	"github.com/mdp/qrterminal/v3"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type savedMessage struct {
	data map[interface{}]string
}

var SavedMessage = savedMessage{
	data: make(map[interface{}]string),
}

func proccessResponse(msg string, sender interface{}, responseChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	if msg == "/reset" {
		SavedMessage.data[sender] = ""
		responseChan <- "Successfully reset the conversation, you can start a new conversation now"
		return
	}

	reqMsg := SavedMessage.data[sender] + "]\n" + msg
	if SavedMessage.data[sender] != "" {
		SavedMessage.data[sender] = SavedMessage.data[sender] + "\nmessage : " + msg
	} else {
		SavedMessage.data[sender] = "\nprevious_messages : [\nmessage : " + msg
		reqMsg = msg
	}

	logrus.Print("\033[31m\n", reqMsg, "\033[0m\n\n")

	response, err := aiclient.AskChatGPT(&reqMsg)
	if err != nil || len(response.Choices) < 1 {
		responseChan <- "Sorry, I don't understand what you mean"
		return
	}

	msg = response.Choices[0].Message.Content
	logrus.Println("\033[34m", msg, "\033[0m")

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

		logrus.Println("\033[32mSender\t:", senderName, " | ", sender, "\033[0m")
		logrus.Println("\033[32mMessage\t:", msg, "\033[0m")
		logrus.Println("\033[32mReply\t:", SavedMessage.data[sender], "\033[0m")

		if debug {
			return
		}

		responseChan := make(chan string, 1)
		var wg sync.WaitGroup
		wg.Add(1)

		go proccessResponse(msg, sender, responseChan, &wg)

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

func SendMessage(client *whatsmeow.Client, clientNumber int, msg string) error {
	protoMsg := &proto.Message{
		ExtendedTextMessage: &proto.ExtendedTextMessage{
			// text to be sent to the sender
			Text: &msg,
		},
	}

	srv := client.Store.ID.Server
	receiver := types.JID{
		User:   fmt.Sprint(clientNumber),
		Server: srv,
	}

	_, err := client.SendMessage(context.Background(), receiver, protoMsg)

	return err
}

func Logout(client *whatsmeow.Client) error {
	err := client.Logout()
	if err != nil {
		logrus.Errorln("Failed to logout. Doing forceful logout", err)
		err = client.Store.Delete()
	}
	return err
}

func GenerateQRCode(client *whatsmeow.Client) error {
	qr, err := client.GetQRChannel(context.Background())
	if err != nil {
		logrus.Errorln("Failed to get QR channel first step", err)
		return err
	}
	logrus.Infoln("QR : ", &qr)
	for evt := range qr {
		if evt.Event == "code" {
			// Render the QR code here
			// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
			fmt.Println("QR code:", evt.Code)

			// render the QR code and wait for the user to scan it
			// then continue
			config := qrterminal.Config{
				Level:     qrterminal.L,
				Writer:    os.Stdout,
				BlackChar: qrterminal.WHITE,
				WhiteChar: qrterminal.BLACK,
				QuietZone: 1,
			}
			qrterminal.GenerateWithConfig(evt.Code, config)
			fmt.Println("Scan the QR code above")
		} else {
			fmt.Println("Login event:", evt.Event)
		}
	}
	return nil
}
