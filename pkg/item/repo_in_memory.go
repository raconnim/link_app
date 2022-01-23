package item

import (
	"errors"
	"fmt"
	"sync"
)

type ItemRepositoryInMemory struct {
	data map[string]string
	mu     *sync.RWMutex
}

var (
	ErrNoLongLink  = errors.New("no such long link")
	ErrNoShortLink = errors.New("no such short link")
)

func NewItemRepositoryInMemory() *ItemRepositoryInMemory {
	repo := make(map[string]string, 10)

	return &ItemRepositoryInMemory{
		data: repo,
		mu: &sync.RWMutex{},
	}
}

func (r *ItemRepositoryInMemory) SearchLongLink(shortLink string) (string, error) {
	r.mu.RLock()
	longLink, ok := r.data[shortLink]
	r.mu.RUnlock()
	if !ok {
		return "", ErrNoShortLink
	}

	return longLink, nil
}

func (r *ItemRepositoryInMemory) searchShortLink(longLink string) (string, error) {
	r.mu.RLock()
	data := r.data
	r.mu.RUnlock()

	for shortLink, currLongLink := range data {
		if longLink == currLongLink {
			return shortLink, nil
		}
	}

	return "", ErrNoLongLink
}

func (r *ItemRepositoryInMemory) addLinkInMemory(longLink, shortLink string) {
	r.mu.Lock()
	r.data[shortLink] = longLink
	r.mu.Unlock()
}

func (r *ItemRepositoryInMemory) checkLink(longLink, shortLink string) error {
	currLink, err := r.SearchLongLink(shortLink)
	if err == ErrNoShortLink || err == nil && currLink == longLink {
		return nil
	}

	return fmt.Errorf("link exist or error with db")
}

func (r *ItemRepositoryInMemory) AddLink(longLink string) (string, error) {
	shortLink, err := r.searchShortLink(longLink)
	if err == nil {
		return shortLink, nil
	}
	//if err != ErrNoLongLink {
	//	return "", err
	//}

	shortLink = GenerateShortLink()

	err = r.checkLink(longLink, shortLink)
	if err != nil {
		return "", err
	}

	r.addLinkInMemory(longLink, shortLink)

	return shortLink, nil
}
