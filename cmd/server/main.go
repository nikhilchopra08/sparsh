package main

import (
	"licensebox/internal/api"
	"licensebox/internal/db"
	"log"
)

func main() {
	// Initialize Database
	database, err := db.InitDB("licenses.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup Router
	r := api.NewRouter(database)
	router := r.SetupRoutes()

	log.Println("🚀 License Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
