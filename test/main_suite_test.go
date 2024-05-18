package main_test

import (
	"log"
	"sync"
	"testing"

	"github.com/farhanaltariq/fiberplate/database"
	"github.com/farhanaltariq/fiberplate/routes"
	"github.com/farhanaltariq/fiberplate/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestMain(t *testing.T) {
	if testExecuted {
		return
	}
	testExecuted = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var (
	app          *fiber.App
	baseUrl      = ":3000"
	fullURL      = "http://127.0.0.1:3000/api"
	once         sync.Once
	testExecuted = false
)

func setupFiberApp() {

	utils.CustomFormatter()

	if err := database.Connect(); err != nil {
		log.Fatal("Error connecting to database", err)
	}

	app = fiber.New(
		fiber.Config{
			Prefork:           false,
			ReduceMemoryUsage: true,
			BodyLimit:         100000 * 1024 * 1024,
			Concurrency:       10000,
		},
	)

	app.Use(recover.New())
	routes.Init(app)

	go func() {
		if err := app.Listen(baseUrl); err != nil {
			logrus.Errorln(err)
		}
	}()

}

var _ = BeforeSuite(func() {
	once.Do(func() {
		setupFiberApp()
	})
})
