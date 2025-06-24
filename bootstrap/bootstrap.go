package bootstrap

import (
	"SimpleMessaging/app/ws"
	"SimpleMessaging/pkg/database"
	"SimpleMessaging/pkg/env"
	"SimpleMessaging/pkg/router"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogger()
	database.SetUpDatabase()
	database.SetupMongoDb()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	go ws.ServeMessaging(app)

	router.InstallRouter(app)

	return app
}

func SetupLogger() {
	fileLogg, err := os.OpenFile("./logs/app_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, fileLogg)
	log.SetOutput(mw)
}
