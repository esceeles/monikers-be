package main

import (
	"encoding/json"
	"fmt"
	"log"
	"monikers/deck"
	"net/http"
)

func main() {
	handleRequests()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
func respond(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
	return
}

func withMessage(message string) interface{} {
	resp := make(map[string]string)
	resp["message"] = message

	return resp
}

func handleRequests() {
	http.HandleFunc("/createGame", createGame)
	http.HandleFunc("/joinGame", joinGame)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

type startGameRequest struct {
	hostName string
}

type joinGameRequest struct {
	name     string
	joinCode string
}

type startGameResponse struct {
	hostHand []int
	joinCode string
}

type joinGameResponse struct {
	hand []int
}

func createGame(w http.ResponseWriter, r *http.Request) {
	//shuffles and stores the deck in the database with the generated idCode and host userName for scoring

	var request startGameRequest

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		respond(w, http.StatusBadRequest, withMessage(fmt.Sprintf("error unmarshalling request:%v", err)))
	}

	hostHand, joinCode, err := deck.CreateGame(request.hostName)
	if err != nil {
		respond(w, http.StatusBadRequest, withMessage(fmt.Sprintf("error creating game: %v", err)))
	}

	response := startGameResponse{
		hostHand: hostHand,
		joinCode: joinCode,
	}

	respond(w, http.StatusCreated, response)
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	//get your cards from the deck and add your username to database for scoring

	var request joinGameRequest

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		respond(w, http.StatusBadRequest, withMessage(fmt.Sprintf("error unmarshalling request:%v", err)))
	}

	hand, err := deck.JoinGame(request.name, request.joinCode)
	if err != nil {
		respond(w, http.StatusBadRequest, withMessage(fmt.Sprintf("error joining game: %v", err)))
	}

	response := joinGameResponse{
		hand: hand,
	}

	respond(w, http.StatusCreated, response)
}
