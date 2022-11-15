package util

import (
	"math/rand"
	"time"
)

var randomList = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q",
	"r", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K",
	"L", "M", "N", "O", "P", "Q", "R", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4",
	"5", "6", "7", "8", "9",
}

func randIdBySize(size int) (ID string) {
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	for i := 0; i != size; i += 1 {
		ID += randomList[randomGenerator.Intn(len(randomList))]
	}

	return
}

func RandId30() (ID string) {
	return randIdBySize(30)
}
