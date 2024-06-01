package cache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func ConnectRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		return nil, errors.New("Missing redis configs.")
	}

	connectionString := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: "",
		DB:       0,
	})

	return redisClient, nil
}
