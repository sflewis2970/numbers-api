package models

import (
	"log"
	"numbers-api/config"
	"numbers-api/messaging"
)

type NumbersModel struct {
	cfgData    *config.CfgData
	redisModel *RedisModel
}

var numbersModel *NumbersModel

func (nm *NumbersModel) AddItem(ID string, nbrData messaging.NumbersAPIData) {
	log.Print("Entering models.AddItem...")

	// Add generated number to request before inserting data
	nbrData.ID = ID

	log.Print("ID: ", nbrData.ID)
	log.Print("message: ", nbrData.Message)
	log.Print("timestamp: ", nbrData.Timestamp)

	insertErr := nm.redisModel.Insert(nbrData)
	if insertErr != nil {
		log.Print("Error inserting item...")
		return
	} else {
		log.Print("item successfully inserted...")
	}

	log.Print("Exiting models.AddItem...")
	log.Print("")
}

func (nm *NumbersModel) FindItem(id string) (messaging.NumbersAPIData, error) {
	log.Print("Entering models.FindItem...")

	nmData, getErr := nm.redisModel.Get(id)
	if getErr != nil {
		log.Print("Error getting item...")
		return messaging.NumbersAPIData{}, getErr
	} else {
		log.Print("item successfully retrieved...")
	}

	log.Print("ID: ", nmData.ID)
	log.Print("Message: ", nmData.Message)
	log.Print("Timestamp: ", nmData.Timestamp)

	log.Print("Exiting models.FindItem...")
	log.Print("")

	return nmData, nil
}

func (ggm *NumbersModel) DeleteItem(id string) error {
	log.Print("Entering models.DeleteItem...")

	nbrData, getErr := ggm.FindItem(id)
	if getErr != nil {
		log.Print("Error getting item...")
		return getErr
	}

	delErr := ggm.DeleteItem(nbrData.ID)
	if delErr != nil {
		log.Print("Error deleting item...")
		return delErr
	}

	log.Print("Exiting models.DeleteItem...")
	log.Print("")

	return nil
}

func NewNumbersModel() *NumbersModel {
	log.Print("Creating model object...")
	guessGameModel := new(NumbersModel)

	// Get config data
	guessGameModel.cfgData = config.NewConfig().LoadCfgData()

	// New model (cacheModel)
	guessGameModel.redisModel = NewRedisModel()

	return guessGameModel
}
