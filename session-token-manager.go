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
	"log"
	"io"
) 

func main() {
	client := redis.NewClient(&redis.Options {
		Addr: "session-token-redis:6379",
		Password: "",
		DB: 0,
	})
	// submit to redis
	err := client.Set("key", "value", 0).Err()
	val, err := client.Get("key").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("key", val)
	http.HandleFunc("/create", createSessionToken)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Message struct {
	Data []byte `json:"data"`
}

type ResponseMessage struct {
	Data string `json:"data"`
}

func parseRequestBody(body io.Reader) ([]byte, error) {
	decoder := json.NewDecoder(body)
	var message Message
	err := decoder.Decode(&message)
	return message.Data, err
}

func writeMessageResponse(w http.ResponseWriter, message ResponseMessage) (error) {
 	w.Header().Set("Content-Type", "application/json")
	messageJson, err := json.Marshal(message)
 	w.Write(messageJson)
	return err
}

func createSessionToken(w http.ResponseWriter, r *http.Request) {
	sessionData, err := parseRequestBody(r.Body)
	print(sessionData)	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	sha1hash := sha1.New()
	sha1hash.Write(sessionData)
	bs := sha1hash.Sum(nil)
	print(bs)
	var message ResponseMessage
	message.Data = string(bs)
	writeMessageResponse(w, message)
}

func deleteSessionToken() {

}

func validateSessionToken() {

}

