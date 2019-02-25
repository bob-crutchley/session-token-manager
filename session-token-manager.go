package main

/*
TODO
create
	accept a request containing session data
	create sha1 hash of session data
	get session data encrypted
	store in redis with expiry time (24hrs for now)
		key: sha1 hash
		value: encrypted session data
delete
	accept post request with session token 
	decrypt data
	create sha1 hash of data
	delete from redis using sha1 as key
validate
	accept post request with session token 
	decrypt data
	create sha1 hash of data
	search redis by sha1 as key
	if the session exists, then valid
*/

import (
	"fmt"
	"github.com/go-redis/redis"
	"crypto/sha1"
	"net/http"
	"encoding/json"
)

func main() {
	client := redis.NewClient(&redis.Options {
		Addr: "session-token-redis:6379",
		Password: "",
		DB: 0,
	})
	sha1hash := sha1.New()
	sha1hash.Write([]byte("string"))
	bs := sha1hash.Sum(nil)
	// submit to redis
	err := client.Set("key", bs, 0).Err()
	val, err := client.Get("key").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("key", val)
}

type Message struct {
	Data []byte `json:"data"`
}

func parseRequestBody(w http.ResponseWriter, body io.Reader) ([]byte, error) {
	decoder := json.NewDecoder(body)
	var message Message
	err := decoder.Decode(&message)
	return message.Data, err
}

func writeMessageResponse(w http.ResponseWriter, message Message) (error) {
 	w.Header().Set("Content-Type", "application/json")
	messageJson, err := json.Marshal(message)
 	w.Write(messageJson)
	return err
}

func create(w http.ResponseWriter, r *http.Request) {
	sessionData := parseRequestBody(w, r.Body)
	if err != nil {
		fmt.Println(err)	
	}

}

