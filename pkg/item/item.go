package item

import (
	"math/rand"
	"time"
)

// change name
type Item struct {
	ShortLink string `json:"short_link,omitempty"`
	LongLink  string `json:"long_link,omitempty"`
}

const letter = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

var GenerateShortLink = func() string {
	shortLink := make([]rune, 10)
	letterRunes := []rune(letter)
	rand.Seed(time.Now().UnixNano())
	for i := range shortLink {
		shortLink[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(shortLink)
}
