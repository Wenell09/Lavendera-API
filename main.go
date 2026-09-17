package main

import (
	"log"

	"github.com/Wenell09/lavendera-api/internal/app"
	"github.com/Wenell09/lavendera-api/internal/shared/config"
)

func main() {
	// Load environment/configuration
	config.LoadConfig()
	// Initialize all dependencies using Google Wire
	application, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	// Close database connection
	sqlDB, err := application.DB.DB()
	if err == nil {
		defer sqlDB.Close()
	}
	// Start server
	log.Printf(
		"Lavendera API running on port %s 🚀",
		config.ENV.Port,
	)
	if err := application.Router.Run(
		":" + config.ENV.Port,
	); err != nil {
		application.Logger.Fatalf(
			"failed to run server: %v",
			err,
		)
	}
}
