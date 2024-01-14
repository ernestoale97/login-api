package storage

import (
	"github.com/go-redis/redis/v7"
	"login_api/pkg/config"
)

var client *redis.Client

func GetRedis() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}
	settings, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	dsn := settings.AppRedis
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client, nil
}
