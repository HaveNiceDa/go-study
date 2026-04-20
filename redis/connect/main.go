package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var DB *redis.Client

func Connect() {
  redisDB := redis.NewClient(&redis.Options{
    Addr:     "127.0.0.1:6379", // 不写默认就是这个
    DB:       0,                // 默认是0
  })
  _, err := redisDB.Ping().Result()
  if err != nil {
    panic(err)
  }
  DB = redisDB
}

func main() {
  Connect()
	fmt.Println(DB.Get("key"))
}