package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eeropu/chatdemo/server"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, world!")
	app := fiber.New()

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		app.Static("/", "./client/dist")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	DB_URI := os.Getenv("DB_URI")
	if DB_URI == "" {
		log.Fatal("DB_URI was empty, can't connect")
	}

	server.Connect(DB_URI)

	defer func() {
		server.Disconnect()
		app.Shutdown()
	}()

	server.Run(app.Group("/api"))

	log.Fatal(app.Listen(":" + PORT))
}
