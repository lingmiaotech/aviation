package tonic

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Redis *redis.Client
var RedisPub *redis.Client
var RedisSub *redis.Client

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

	pubEnabled := Configs.GetBool("redis.pub_enabled")
	if pubEnabled {
		RedisPub = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: "",
			DB:       db,
		})
	}

	subEnabled := Configs.GetBool("redis.sub_enabled")
	if subEnabled {
		RedisSub = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: "",
			DB:       db,
		})
	}

	return
}
