package main

import (
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/samwarnick/perfect-stack-go/db"
	"github.com/samwarnick/perfect-stack-go/models"
	"github.com/samwarnick/perfect-stack-go/pages"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()

	app := fiber.New()

	app.Use(logger.New())

	app.Static("/assets", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		return renderIndex(c)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		var payload models.CreateMessageSchema

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}

		// err = validate.Struct(payload)
		// if err != nil {
		// 	// fmt.Print(err)
		// 	return c.Status(fiber.StatusBadRequest).JSON(err)
		// }

		now := time.Now()
		newMessage := models.Message{
			CreatedAt: now,
			Message:   payload.Message,
		}

		result := db.DB.Create(&newMessage)

		if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist, please use another title"})
		} else if result.Error != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
		}
		return renderIndex(c)
	})

	log.Fatal(app.Listen(":3000"))
}

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}

func renderIndex(c *fiber.Ctx) error {
	name := os.Getenv("NAME")
	var messages []models.Message
	db.DB.Find(&messages)
	return Render(c, pages.Index(name, messages))
}
