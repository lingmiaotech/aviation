package tonic

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func InitRedis() (err error) {

	enabled := Configs.GetBool("redis.enabled")
	if enabled {
		return nil
	}

	host := Configs.GetString("redis.host")
	port := Configs.GetInt("redis.port")
	db := Configs.GetInt("redis.db")

	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return
}

func GetPub() *redis.Client {
	host := Configs.GetString("redis.host")
	port := Configs.GetInt("redis.port")
	db := Configs.GetInt("redis.db")

	RedisPub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return RedisPub
}

func GetSub() *redis.Client {
	host := Configs.GetString("redis.host")
	port := Configs.GetInt("redis.port")
	db := Configs.GetInt("redis.db")

	RedisSub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	return RedisSub
}
