package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
	}
}

func (s *RedisService) Set(value interface{}, key string) error {
	return s.Client.Set(context.TODO(), key, value, time.Hour*24*7).Err()
}

func (s *RedisService) Get(key string) (string, error) {
	result, err := s.Client.Get(context.TODO(), key).Result()
	if err != nil && err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (s *RedisService) Delete(key string) error {
	return s.Client.Del(context.TODO(), key).Err()
}
