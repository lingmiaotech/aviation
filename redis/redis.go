package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lingmiaotech/tonic/configs"
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

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return
}

func GetPub() *redis.Client {
	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")

	PubClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return PubClient
}

func GetSub() *redis.Client {
	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")

	SubClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return SubClient
}
