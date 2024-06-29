package models

import (
	"context"
	"encoding/json"
	"log"
	"numbers-api/config"
	"numbers-api/messaging"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// REDIS_TLS_URL Redis Constants
	REDIS_PASSWORD         string = "REDIS_PASSWORD"
	REDIS_DB_NAME_MSG      string = "GO_REDIS: "
	REDIS_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	REDIS_MARSHAL_ERROR        string = "Marshaling error...: "
	REDIS_UNMARSHAL_ERROR      string = "Unmarshalling error...: "
	REDIS_INSERT_ERROR         string = "Insert error...: "
	REDIS_ITEM_NOT_FOUND_ERROR string = "Item not found...: "
	REDIS_GET_ERROR            string = "Get error...: "
	REDIS_DELETE_ERROR         string = "Delete error...: "
	REDIS_PING_ERROR           string = "Error pinging in-memory cache server...: "
)

type Redis struct {
	TLS_URL  string `json:"tls_url"`
	URL      string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type RedisModel struct {
	cfgData  *config.CfgData
	memCache *redis.Client
}

var redisModel *RedisModel

// Ping database server, since this is local to the server make sure the object for storing data is created
func (rm *RedisModel) Ping() error {
	log.Print("Entering models.Ping...")

	ctx := context.Background()

	statusCmd := rm.memCache.Ping(ctx)
	pingErr := statusCmd.Err()
	if pingErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_PING_ERROR, pingErr)
		return pingErr
	}

	log.Print("Exiting models.Ping...")
	log.Print("")

	return nil
}

// Insert a single record into table
func (rm *RedisModel) Insert(nData messaging.NumbersAPIData) error {
	log.Print("Entering models.Insert...")

	ctx := context.Background()

	byteStream, marshalErr := json.Marshal(nData)
	if marshalErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_MARSHAL_ERROR, marshalErr)
		return marshalErr
	}

	setErr := rm.memCache.Set(ctx, nData.ID, byteStream, time.Duration(0)).Err()
	if setErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_INSERT_ERROR, setErr)
		return setErr
	}

	log.Print("data inserted...")
	log.Print("Exiting modelsInsert...")

	return nil
}

// Get a single record from table
func (rm *RedisModel) Get(id string) (messaging.NumbersAPIData, error) {
	log.Print("Entering models.Get...")

	ctx := context.Background()

	log.Print("Getting record from the map, with ID: ", id)

	var nData messaging.NumbersAPIData
	getResult, getErr := rm.memCache.Get(ctx, id).Result()
	if getErr == redis.Nil {
		log.Print(REDIS_DB_NAME_MSG + REDIS_ITEM_NOT_FOUND_ERROR)
		return messaging.NumbersAPIData{}, nil
	} else if getErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_GET_ERROR, getErr)
		return messaging.NumbersAPIData{}, getErr
	} else {
		unmarshalErr := json.Unmarshal([]byte(getResult), &nData)
		if unmarshalErr != nil {
			log.Print(REDIS_DB_NAME_MSG+REDIS_UNMARSHAL_ERROR, unmarshalErr)
			return messaging.NumbersAPIData{}, unmarshalErr
		}

		log.Print("data retrieved...")
	}

	log.Print("Exiting models.Get...")
	log.Print("")

	return nData, nil
}

// Update a single record in table
func (rm *RedisModel) Update(nUpdateData messaging.NumbersAPIData) error {
	log.Print("Entering models.Update...")

	ctx := context.Background()

	log.Println("Updating record in the map")

	// Let's make sure that the item is already in the DB
	_, nErr := rm.Get(nUpdateData.ID)
	if nErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_GET_ERROR, nErr)
		return nErr
	}

	// Convert the data to a byte stream
	byteStream, marshalErr := json.Marshal(nUpdateData)
	if marshalErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_MARSHAL_ERROR, marshalErr)
		return marshalErr
	}

	// Update the data
	updateErr := rm.memCache.Set(ctx, nUpdateData.ID, byteStream, time.Duration(0)).Err()
	if updateErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_INSERT_ERROR, updateErr)
		return updateErr
	}

	log.Print("data updated...")
	log.Print("Exiting models.Update...")
	log.Print("")

	return nil
}

// Delete a single record from table
func (rm *RedisModel) Delete(id string) error {
	log.Print("Entering models.Delete...")
	log.Print("Deleting record with ID: ", id)

	// Delete the record from map
	ctx := context.Background()
	delErr := rm.memCache.Del(ctx, id).Err()
	if delErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_DELETE_ERROR, delErr)
		return delErr
	}

	log.Print("Exiting models.Delete...")
	log.Print("")

	return nil
}

func NewRedisModel() *RedisModel {
	// Initialize go-cache in-memory cache model
	log.Print("Creating goRedis dbModel object...")
	redisModel = new(RedisModel)

	// Get config data
	redisModel.cfgData = config.NewConfig().LoadCfgData()

	// Define go-redis cache settings
	log.Print(REDIS_DB_NAME_MSG + REDIS_CREATE_CACHE_MSG)

	// Define connection variables
	var redisOptions *redis.Options

	// The config package handles reading the environment variables and parsing the url.
	// Once the external packages access the values, the environment has already been taken
	// care of.
	redisAddr := redisModel.cfgData.Redis.URL + ":" + redisModel.cfgData.Redis.Port
	log.Print("The redis address is...: ", redisAddr)

	redisOptions = &redis.Options{
		Addr:     redisAddr, // redis Server Address,
		Password: "",        // set password
		DB:       0,         // use default DB
	}

	// Create go-redis in-memory cache
	redisModel.memCache = redis.NewClient(redisOptions)

	return redisModel
}
