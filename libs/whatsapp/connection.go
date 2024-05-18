package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/farhanaltariq/fiberplate/utils"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func ConnectDB(client *whatsmeow.Client) error {
	var err error
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
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
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func Init() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	psqlInfo := utils.GetEnv("CRED_DB_DSN", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	store, err := sqlstore.New("postgres", psqlInfo, dbLog)
	if err != nil {
		panic(err)
	}

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := store.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	// defer client.Disconnect()

	// Use a channel to receive signals for graceful shutdown
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the database connection concurrently
	go func() {
		err = ConnectDB(client)
		if err != nil {
			panic(err)
		}
	}()

	// Start the event handling concurrently
	go func() {
		client.AddEventHandler(func(evt interface{}) {
			// set to true to not send result to client
			EventHandler(client, evt, false)
		})
	}()
}
