package main

import (
	"net/http"
	"time"
	controllers "undetectable-ai/DSk/controller"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

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
		"http://127.0.0.1:5500":                 true,
		"https://humanize-ai-server.vercel.app": true,
	}

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
		AllowMethods:     "POST, OPTIONS",
		AllowHeaders:     "Content-Type, x-api-key",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Server is running on Vercel!")
	})

	app.Post("/rewrite", controllers.RewriteHandler)

	return app
}

func VercelHandler(w http.ResponseWriter, r *http.Request) {
	adaptor.FiberApp(setupApp())(w, r)
}

// I commented this for production sake
// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	log.Printf("🚀 Server running on port %s\n", port)
// 	http.HandleFunc("/", VercelHandler)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }
