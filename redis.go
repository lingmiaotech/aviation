package tonic

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lingmiaotech/tonic/configs"
)

var Redis *redis.Client

func InitRedis() (err error) {

	enabled := configs.GetBool("redis.enabled")
	if enabled {
		return nil
	}

	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")

	Redis = redis.NewClient(&redis.Options{
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

	RedisPub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return RedisPub
}

func GetSub() *redis.Client {
	host := configs.GetString("redis.host")
	port := configs.GetInt("redis.port")
	db := configs.GetInt("redis.db")

	RedisSub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return RedisSub
}
