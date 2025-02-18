package handler

import (
	"log"
	"net/http"
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
	app = setupApp()
	log.Println("🚀 Initializing server...")
}

func setupApp() *fiber.App {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	allowedOrigins := map[string]bool{
		"http://127.0.0.1:5500":                   true,
		"https://humanize-ai-frontend.vercel.app": true, // NODE
		"https://humanize-ai-one.vercel.app":      true, // GO
		"https://humanize-ai-server.vercel.app":   true, // PY
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
		return c.SendString("🚀 Server is running on Vercel!")
	})

	app.Post("/rewrite", controllers.RewriteHandler)

	return app
}

func Handler(w http.ResponseWriter, r *http.Request) {
	bridge.Start(adaptor.FiberApp(app))
}
