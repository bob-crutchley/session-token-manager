package main

/*
TODO
create
	store in redis with expiry time (24hrs for now)
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
	"bytes"
) 
var redisClient *redis.Client
func main() {
	redisClient = redis.NewClient(&redis.Options {
		Addr: "session-token-redis:6379",
		Password: "",
		DB: 0,
	})
	// submit to redis
	err := redisClient.Set("key", "value", 0).Err()
	val, err := redisClient.Get("key").Result()
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

func parseRequestBody(body io.Reader) (Message, error) {
	decoder := json.NewDecoder(body)
	var message Message
	err := decoder.Decode(&message)
	return message, err
}

func writeMessageResponse(w http.ResponseWriter, message Message) (error) {
 	w.Header().Set("Content-Type", "application/json")
	messageJson, err := json.Marshal(message)
 	w.Write(messageJson)
	return err
}

func createSessionToken(w http.ResponseWriter, r *http.Request) {
	session, err := parseRequestBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(session.Data)
	sessionBytesRepresentation, err := json.Marshal(session)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.Post("http://aes-crypto:8000/encrypt", "application/json", bytes.NewBuffer(sessionBytesRepresentation))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var encryptedMessage Message
	json.NewDecoder(resp.Body).Decode(&encryptedMessage)
	fmt.Println(encryptedMessage.Data)

	sha1hash := sha1.New()
	sha1hash.Write(session.Data)
	bs := sha1hash.Sum(nil)
	fmt.Println(bs)

	err = redisClient.Set(string(bs), encryptedMessage.Data, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	val, err := redisClient.Get(string(bs)).Result()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println(bs, val)

	var message Message
	message.Data = bs
	writeMessageResponse(w, message)
}

func deleteSessionToken() {

}

func validateSessionToken() {

}

