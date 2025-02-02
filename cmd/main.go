package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdulmalikraji/verify-influencers-backend/config"
	"github.com/abdulmalikraji/verify-influencers-backend/db/connection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {

	//load ennvironment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	client := connection.New()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(recover.New())

	config.InitializeRoutes(app, client)

	// testCase, err := twitter.GetTwitterUserByUsername("drgabriellelyon")
	// if err != nil {
	// 	log.Fatalf("Error finding user: %v\n", err)
	// }
	// fmt.Println("Test result: ", *testCase)

	// Start the server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "4000" // Default for local dev
		}
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Call gracefulShutdown to handle cleanup
	gracefulShutdown(app, client)
}

func gracefulShutdown(app *fiber.App, client connection.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	// Close the PostgreSQL database connection
	database, err := client.PostgresConnection.DB()
	if err != nil {
		log.Println("PostgreSQL Closing ERROR :", err)
	}
	database.Close()
	log.Printf("PostgreSQL Closed")

	// Attempt to shut down the Fiber app gracefully
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v\n", err)
	}

	log.Println("Server shutdown complete.")
}
