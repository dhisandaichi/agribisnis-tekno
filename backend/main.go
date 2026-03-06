package main

import (
	"log"

	"github.com/agritec/backend/config"
	"github.com/agritec/backend/routes"
)

func main() {
	// Initialize configurations (DB, Env config)
	config.InitConfig()
	db := config.ConnectDB()

	// Initialize Routes
	router := routes.SetupRouter(db)

	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
