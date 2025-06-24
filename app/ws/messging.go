package ws

import (
	"SimpleMessaging/app/models"
	"SimpleMessaging/app/repository"
	"SimpleMessaging/pkg/env"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServeMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan models.MessagePayload)

	app.Get("message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true

		for {
			var m models.MessagePayload
			if err := c.ReadJSON(&m); err != nil {
				log.Println("ERROR PAYLOAD: ", err)
				break
			}
			m.Date = time.Now()
			err := repository.InsertNewMessage(context.Background(), m)
			if err != nil {
				log.Println(err)
			}
			broadcast <- m
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for c := range clients {
				err := c.WriteJSON(msg)
				if err != nil {
					log.Println("failed to write json: ", err)
					c.Close()
					delete(clients, c)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT_SOCKET", "8080"))))
}
