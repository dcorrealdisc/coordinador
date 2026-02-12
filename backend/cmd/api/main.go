package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Crear instancia de Fiber
	app := fiber.New(fiber.Config{
		AppName: "Coordinador API v0.1.0",
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "coordinador-api",
			"version": "0.1.0",
		})
	})

	// Rutas API
	api := app.Group("/api/v1")

	// Placeholder endpoints para cada módulo
	api.Get("/students", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Students endpoint - Coming soon"})
	})

	api.Get("/courses", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Courses endpoint - Coming soon"})
	})

	api.Get("/planning", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Planning endpoint - Coming soon"})
	})

	api.Get("/reports", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Reports endpoint - Coming soon"})
	})

	api.Get("/tutors", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Tutors endpoint - Coming soon"})
	})

	// Iniciar servidor
	port := ":8080"
	log.Printf("🚀 Coordinador API iniciando en http://localhost%s", port)
	log.Fatal(app.Listen(port))
}
