package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

)

func main() {
	fmt.Println("Starting server...")

	app := fiber.New()

	app.Get("/abc", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	log.Fatal(app.Listen(":3000"))
}
