package main

import (
	"log"
	"net/http"
	"os"
	"time"
	controllers "undetectable-ai/DSk/controller"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/vercel/go-bridge/go/bridge"
)

var app *fiber.App

func init() {
	// Initialize Fiber app once
	app = fiber.New()

	// Common middleware setup
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	allowedOrigins := map[string]bool{
		"http://127.0.0.1:5500":                   true,
		"https://humanize-ai-frontend.vercel.app": true,
		"https://humanize-ai-one.vercel.app":      true,
		"https://humanize-ai-server.vercel.app":   true,
		"http://localhost:3000":                   true,
		"http://localhost:8080":                   true,
	}

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
		AllowMethods:     "POST",
		AllowHeaders:     "Content-Type, X-API-Key",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Combined Server is running!")
	})

	app.Post("/rewrite", controllers.RewriteHandler)
}

func StartServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("🚀 Server running on port %s", port)
	return app.Listen(":" + port)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	bridge.Start(adaptor.FiberApp(app))
}

func main() {
	if err := StartServer(); err != nil {
		log.Fatal(err)
	}
}
