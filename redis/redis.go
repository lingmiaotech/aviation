package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/dyliu/tonic/configs"
)

var Client *redis.Client

func InitRedis() (err error) {

	enabled := configs.GetBool("redis.enabled")
	if !enabled {
		return nil
	}

	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")
	password := configs.GetDynamicString("redis.password")

	Client = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		Password:   password,
		DB:         db,
		MaxRetries: 3,
	})

	return
}

func GetPub() *redis.Client {
	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")
	password := configs.GetDynamicString("redis.password")

	PubClient := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		Password:   password,
		DB:         db,
		MaxRetries: 3,
	})

	return PubClient
}

func GetSub() *redis.Client {
	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")
	password := configs.GetDynamicString("redis.password")

	SubClient := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		Password:   password,
		DB:         db,
		MaxRetries: 3,
	})

	return SubClient
}
