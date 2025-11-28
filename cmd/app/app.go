package app

import (
	"log"

	"aadhaar-user-service/internals/config"
	"aadhaar-user-service/internals/database"
	"aadhaar-user-service/internals/server"
)

func Setup() {
	database.Connect()

	config.Automigration()

	server.Setup()
	app := server.New()

	log.Println("Starting Aadhaar User Service on port :3015")
	if err := app.Listen(":3015"); err != nil {
		log.Fatalf("Error starting server %v\n", err)
	}
}
