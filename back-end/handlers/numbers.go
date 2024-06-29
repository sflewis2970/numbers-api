package handlers

import (
	"log"
	"net/http"
	"numbers-api/messaging"
	"numbers-api/models"
	"time"
)

type NumbersHandler struct {
	numbersModel *models.NumbersModel
}

var nbrsHandler *NumbersHandler

func (nbrsh *NumbersHandler) StartNumbersService(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.StartNumbersService...")

}

func (nh *NumbersHandler) GetDateFact(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.GetDateFact...")

	// Get game ID from query parameter
	nbrsID := r.URL.Query().Get("id")

	// Get guessed number from query parameter
	dateStr := r.URL.Query().Get("date")

	log.Print("Query Parameter, id...:", nbrsID)
	log.Print("Query Parameter, guess...:", dateStr)

	// Set content-type
	rw.Header().Add("Content-Type", "application/json")

	var nbrsData messaging.NumbersAPIData

	_, convErr := time.Parse(dateStr, "YYYY-MM-DD")
	if convErr != nil {
		log.Print("Error converting string to int...: ", convErr)

		// Update AnswerResponse
		nbrsData.Message = convErr.Error()

		// Update HTTP Header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		return
	}

	// Check guessed number
	var getErr error

	log.Print("Sending gameID to model...")
	nbrsData, getErr = nh.numbersModel.FindItem(nbrsID)
	if getErr != nil {
		log.Print("Error converting string to int...: ", convErr)
		// Update AnswerResponse
		nbrsData.Message = getErr.Error()

		// Update HTTP Header
		rw.WriteHeader(http.StatusInternalServerError)

		// Write JSON to stream
		return
	}

	log.Print("Numbers ID: ", nbrsID)
	log.Print("Message: ", nbrsData.Message)
	log.Print("time stamp: ", nbrsData.Timestamp)

	// Send OK status
	rw.WriteHeader(http.StatusOK)

	// Display a log message
	log.Print("Sending response to client...")
}

func (nh *NumbersHandler) GetMathFact(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.GetMathFact...")

}

func (nh *NumbersHandler) GetRandomFact(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.GetRandomFact...")

}

func (nh *NumbersHandler) GetTriviaFact(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.GetRandomFact...")

}

func (nh *NumbersHandler) GetYearFact(rw http.ResponseWriter, r *http.Request) {
	log.Print("Entering handlers.GetRandomFact...")

}

func NewNumbersHandler() *NumbersHandler {
	nbrsHandler = new(NumbersHandler)

	nbrsHandler.numbersModel = models.NewNumbersModel()

	return nbrsHandler
}
