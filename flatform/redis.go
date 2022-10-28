package flatform

import (
	"errors"

	"github.com/go-redis/redis"
)

func NewInitResdis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if redisClient == nil {
		return nil, errors.New("can not connect redis")
	}

	return redisClient, nil
}
