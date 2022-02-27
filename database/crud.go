package database

import "log"

func CreateGame(deck []int, joinCode string, hostName string) error {
	log.Println("created game")
	return nil
}

func GetDeck(joinCode string) ([]int, error) {
	log.Println("got deck")
	return nil, nil
}

func JoinGame(deck []int, joinCode string, name string) error {
	log.Println("joined game")
	return nil
}
