package config

import (
	"log"

	"github.com/go-redis/redis"
)

// RC Context
var RC *redis.Client

// Redis Initialization
func Redis() {

	RC = redis.NewClient(&redis.Options{Addr: LoadConfig().CacheURL, Password: LoadConfig().CachePassword, DB: 0})

	_, err := RC.Ping().Result()

	if err != nil {
		log.Panic(err)
	}
	log.Println("Redis Connected")
}
