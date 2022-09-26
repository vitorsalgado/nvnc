package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	srv := fiber.New(fiber.Config{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ext := make(chan struct{}, 1)

	go func() {
		<-c

		err := srv.Shutdown()
		if err != nil {
			log.Fatal(fmt.Errorf("server shutdown failed. reason=%w", err))
			return
		}

		ext <- struct{}{}
	}()

	srv.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	if err := srv.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

	<-ext

	fmt.Print("server closed")
}
