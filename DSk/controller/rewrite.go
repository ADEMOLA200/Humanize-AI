package controllers

import (
	"log"
	"undetectable-ai/DSk/services"

	"github.com/gofiber/fiber/v2"
)

type RewriteRequest struct {
	Text string `json:"text" validate:"required,min=10"`
}

type RewriteResponse struct {
	RewrittenText string `json:"rewritten_text"`
	Success       bool   `json:"success"`
}

func RewriteHandler(c *fiber.Ctx) error {
	log.Printf("Incoming request headers: %+v", c.GetReqHeaders())

	type Request struct {
		Text string `json:"text"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Body parse error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Text == "" {
		log.Println("Empty text received")
		return c.Status(400).JSON(fiber.Map{"error": "Text cannot be empty"})
	}

	rewritten := services.RewriteText(req.Text)

	response := RewriteResponse{
		RewrittenText: rewritten,
		Success:       true,
	}

	return c.JSON(response)
}
