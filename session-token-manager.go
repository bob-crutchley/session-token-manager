package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"crypto/sha1"
)

func main() {
	client := redis.NewClient(&redis.Options {
		Addr: "session-token-redis:6379",
		Password: "",
		DB: 0,
	})
	pong, err := client.Ping().Result()
	test := sha1.Sum([]byte("string"))
	fmt.Println(string(test[:]))
	fmt.Println(pong, err)
	err = client.Set("key", "value", 0).Err()
	val, err := client.Get("key").Result()
	fmt.Println("key", val)
}

