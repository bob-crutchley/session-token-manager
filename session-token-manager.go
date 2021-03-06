package main

/*
TODO
delete
	accept post request with session token 
	decrypt data
	create sha1 hash of data
	delete from redis using sha1 as key
*/

import (
	"github.com/go-redis/redis"
	"crypto/sha1"
	"net/http"
	"encoding/json"
	"log"
	"io"
	"bytes"
	"time"
) 

var redisClient *redis.Client

func main() {
	redisClient = getRedisClient()
	http.HandleFunc("/create", createSessionToken)
	http.HandleFunc("/getsession", getSession)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getRedisClient() (*redis.Client) {
	return redis.NewClient(&redis.Options {
		Addr: "session-token-redis:6379",
		Password: "",
		DB: 0,
	})
}

type Message struct {
	Data []byte `json:"data"`
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

func sessionTokenCreationFailure(w http.ResponseWriter) {
		http.Error(w, "failed to create session token", 500)
}

func getSha1Hash(data []byte) ([]byte) {
	sha1hash := sha1.New()
	sha1hash.Write(data)
	return sha1hash.Sum(nil)
}

func createSessionToken(w http.ResponseWriter, r *http.Request) {
	newSession, err := parseRequestBody(r.Body)
	if err != nil {
		sessionTokenCreationFailure(w)
		return
	}
	newSessionBytesRepresentation, err := json.Marshal(newSession)
	if err != nil {
		sessionTokenCreationFailure(w)
		return
	}
	resp, err := http.Post("http://aes-crypto:8000/encrypt", "application/json", bytes.NewBuffer(newSessionBytesRepresentation))
	if err != nil {
		sessionTokenCreationFailure(w)
		return
	}
	var encryptedMessage Message
	json.NewDecoder(resp.Body).Decode(&encryptedMessage)
	sessionTokenKey := getSha1Hash(encryptedMessage.Data)
	err = redisClient.Set(string(sessionTokenKey), encryptedMessage.Data, 24 * time.Hour).Err()
	if err != nil {
		sessionTokenCreationFailure(w)
		return
	}
	var message Message
	message.Data = sessionTokenKey
	writeMessageResponse(w, message)
}

func deleteSessionToken() {

}

func sessionRetreivalFailure(w http.ResponseWriter) {
		http.Error(w, "failed to retrieve session token", 500)
}

func getSession(w http.ResponseWriter, r *http.Request) {
	sessionRequest, err := parseRequestBody(r.Body)
	if err != nil {
		sessionRetreivalFailure(w)
		return
	}
	sessionRequestKey := sessionRequest.Data

	encryptedSessionData, err := redisClient.Get(string(sessionRequestKey)).Result()
	if err != nil {
		sessionRetreivalFailure(w)
		fmt.Println(err)
	}

	var encryptedSession Message
	encryptedSession.Data = []byte(encryptedSessionData)

	encryptedSessionBytesRepresentation, err := json.Marshal(encryptedSession)
	if err != nil {
		sessionRetreivalFailure(w)
		return
	}

	resp, err := http.Post("http://aes-crypto:8000/decrypt", "application/json", bytes.NewBuffer(encryptedSessionBytesRepresentation))
	if err != nil {
		sessionRetreivalFailure(w)
		return
	}

	var session Message
	json.NewDecoder(resp.Body).Decode(&session)
	writeMessageResponse(w, session)
}

