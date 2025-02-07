package main

import (
	"log"
	"time"
	controllers "undetectable-ai/DSk/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:5500",
		AllowMethods:     "POST, OPTIONS",
		AllowHeaders:     "Content-Type",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

	app.Post("/rewrite", controllers.RewriteHandler)

	log.Println("Server started at :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
