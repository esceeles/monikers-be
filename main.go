package main

import (
	"encoding/json"
	"fmt"
	"log"
	"monikers/deck"
	"net/http"
)

func main() {
	http.HandleFunc("/createGame", createGame)
	http.HandleFunc("/joinGame", joinGame)
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

type startGameResponse struct {
	hostHand []int
	joinCode string
}

type joinGameResponse struct {
	hand []int
}

func createGame(w http.ResponseWriter, r *http.Request) {
	//shuffles and stores the deck in the database with the generated idCode and host userName for scoring

	hostName, ok := r.URL.Query()["hostName"]
	if !ok || len(hostName[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	hostHand, joinCode, err := deck.CreateGame(hostName[0])
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

	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	joinCode, ok := r.URL.Query()["joinCode"]
	if !ok || len(joinCode[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	
	hand, err := deck.JoinGame(name[0], joinCode[0])
	if err != nil {
		respond(w, http.StatusBadRequest, withMessage(fmt.Sprintf("error joining game: %v", err)))
	}

	response := joinGameResponse{
		hand: hand,
	}

	respond(w, http.StatusCreated, response)
}
