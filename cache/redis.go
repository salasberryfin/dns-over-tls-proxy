package cache

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

const (
	redisPort = "6379"
	redisHost = "redis-server"
)

// Connect is a naive function that connects to a Redis server
func Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})
	if client == nil {
		log.Fatalf("Failed to configure Redis client\n")
	}

	msg, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("No connection to Redis: %v\n", err)
	}
	log.Printf("Message from Redis %v\n", msg)
}
