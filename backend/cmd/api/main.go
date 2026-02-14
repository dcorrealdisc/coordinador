package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/dcorreal/coordinador/internal/database"
	"github.com/dcorreal/coordinador/internal/handlers"
	"github.com/dcorreal/coordinador/internal/repositories"
	"github.com/dcorreal/coordinador/internal/services"
)

func main() {
	// Database connection
	dbConfig := database.ConfigFromEnv()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Dependency injection: Repository -> Service -> Handler
	studentRepo := repositories.NewStudentRepository(db)
	catalogRepo := repositories.NewCatalogRepository(db)
	studentService := services.NewStudentService(studentRepo)
	studentImportService := services.NewStudentImportService(studentService, studentRepo, catalogRepo)
	studentHandler := handlers.NewStudentHandler(studentService, studentImportService)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Coordinador API v0.1.0",
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		if err := database.HealthCheck(db); err != nil {
			dbStatus = "error"
		}
		return c.JSON(fiber.Map{
			"status":   "ok",
			"service":  "coordinador-api",
			"version":  "0.1.0",
			"database": dbStatus,
		})
	})

	// API routes
	api := app.Group("/api/v1")
	studentHandler.RegisterRoutes(api)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Coordinador API starting on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
