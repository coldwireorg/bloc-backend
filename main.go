package main

import (
	"log"
	"os"

	"bloc/database"
	"bloc/routes"
	"bloc/utils"
	"bloc/utils/tokens"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// Load env values
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Unable to load env values: %v\n", err)
	} else {
		log.Println("Loaded env values successfully")
	}

	// Check if storage path exist
	utils.CheckPaths()
}

func main() {
	// Create fiber instance
	app := fiber.New(fiber.Config{
		BodyLimit: (1024 * 1024 * 1024) * 8, // Limit of files size (8Gb)
	})

	app.Use(cors.New()) // Add cors

	// Generate ed25519 keypair for signing JWT tokens
	// TODO: allowing to put our own private/public key in env variable,
	// 			 so then we can scale Bloc without worrying about authentication issues
	tokens.GenerateKeys()

	// Create the connection to the database database
	err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database!!", err)
		panic(err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	routes.SetupRoutes(app) // Import routes

	// Get IP and Port to listen on
	listner := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")

	// Create listener
	err = app.Listen(listner)
	if err != nil {
		log.Println("Unable to start server on [" + listner + "]")
		panic(err)
	} else {
		log.Println("Listening on [" + listner + "]")
	}
}
