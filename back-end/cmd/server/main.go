package main

import (
	"log"
	"net/http"
	"numbers-api/config"
	"numbers-api/controllers"

	"github.com/rs/cors"
)

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load config data
	cfgData := config.NewConfig().LoadCfgData()

	// Create controller
	controller := controllers.NewController()

	// setup Cors
	log.Print("Setting up CORS...")
	corsOptionsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	corsHandler := corsOptionsHandler.Handler(controller.Router)

	// Server Address info
	addr := cfgData.Host + ":" + cfgData.Port
	log.Print("The address used by the service is: ", addr)

	// Start Server
	log.Print("Web server is ready...")

	// Listen and Serve
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}
