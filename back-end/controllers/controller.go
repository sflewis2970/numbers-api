package controllers

import (
	"log"
	"numbers-api/handlers"

	"github.com/gorilla/mux"
)

// Controller structure defines teh layout of the Controller
type Controller struct {
	Router         *mux.Router
	NumbersHandler *handlers.NumbersHandler
}

// Package controllers object
var controller *Controller

func (c *Controller) setupRoutes() {
	// Display log message
	log.Print("Setting up service routes")

	// API routes
	c.Router.HandleFunc("/api/v1/numbers/startnbrsservice", c.NumbersHandler.StartNumbersService).Methods("GET")
	c.Router.HandleFunc("/api/v1/numbers/datefact", c.NumbersHandler.GetDateFact).Methods("GET")
	c.Router.HandleFunc("/api/v1/numbers/mathfact", c.NumbersHandler.GetMathFact).Methods("GET")
	c.Router.HandleFunc("/api/v1/numbers/randomfact", c.NumbersHandler.GetRandomFact).Methods("GET")
	c.Router.HandleFunc("/api/v1/numbers/triviafact", c.NumbersHandler.GetTriviaFact).Methods("GET")
	c.Router.HandleFunc("/api/v1/numbers/yearfact", c.NumbersHandler.GetYearFact).Methods("GET")

}

// NewController function create a new Controller and initializes new Controller object
func NewController() *Controller {
	if controller != nil {
		log.Print("Returning controllers object...")
		return controller
	}

	// Create controllers component
	log.Print("Creating controllers object...")
	controller = new(Controller)

	// Numbers handler
	controller.NumbersHandler = handlers.NewNumbersHandler()

	// Set controllers routes
	controller.Router = mux.NewRouter()
	controller.setupRoutes()

	return controller
}
