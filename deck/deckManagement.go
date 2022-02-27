package deck

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"monikers/database"
	"time"
)

const imageDir = "/Users/elliesceeles/monikers/images"
const joinCodeLength = 8
const numDrawCards = 8

func CreateGame(hostName string) ([]int, string, error) {
	deck := getShuffledDeck()
	deck, hostHand := drawCards(deck)
	joinCode := createGameCode()

	err := database.CreateGame(deck, joinCode, hostName)
	if err != nil {
		return nil, "", fmt.Errorf("error storing create game: %v", err)
	}

	return hostHand, joinCode, nil
}

func JoinGame(name string, joinCode string) ([]int, error) {
	deck, err := database.GetDeck(joinCode)
	deck, hand := drawCards(deck)

	err = database.JoinGame(deck, joinCode, name)
	if err != nil {
		return nil, fmt.Errorf("error storing join game: %v", err)
	}

	return hand, nil
}

func createGameCode() string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, joinCodeLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func drawCards(deck []int) ([]int, []int) {
	var card int

	newHand := make([]int, 0)
	for i := 0; i < numDrawCards; i++ {
		deck, card = draw(deck)

		newHand = append(newHand, card)
	}

	return deck, newHand
}

func draw(deck []int) ([]int, int) {
	card := deck[0]

	return append(deck[:0], deck[1:]...), card
}

func getShuffledDeck() []int {
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	deck := make([]int, len(files))
	for i := 0; i < 10; i++ {
		deck = append(deck, i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	return deck
}
