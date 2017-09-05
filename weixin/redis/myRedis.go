package myRedis

import "github.com/go-redis/redis"

import "../log"

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func init() {
	_, errRedis := RedisClient.Ping().Result()
	if errRedis != nil {
		myLog.Fatal.Fatalf("redis error: %s ", errRedis)
	}
}
