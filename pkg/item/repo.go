package item

import (
	"database/sql"
	"fmt"
)

type ItemRepostory struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepostory {
	return &ItemRepostory{DB: db}
}


// search shortLink in db
func(r *ItemRepostory) SearchLongLink(shortLink string) (string, error) {
	longLink := ""

	err := r.DB.QueryRow(`SELECT long_link FROM dbname WHERE short_link = $1`,
		shortLink).
		Scan(&longLink)
	if err != nil {
		return "", err // TODO change returning error
	}

	return longLink, nil
}

func (r *ItemRepostory) searchShortLink(longLink string) (string, error){
	shortLink := ""
	err := r.DB.QueryRow(`SELECT short_link FROM dbname WHERE long_link = $1`, longLink).
		Scan(&shortLink)
	if err != nil {
		return "", err // TODO change returning error
	}

	return shortLink, nil
}


func (r *ItemRepostory) checkLink(longLink, shortLink string) error {
	currLink, err := r.SearchLongLink(shortLink)
	if err == sql.ErrNoRows || err == nil && currLink == longLink {
		return nil
	}

	return fmt.Errorf("link exist or error with db")
}

func (r *ItemRepostory) addLinkInDB(longLink, shortLink string) error {
	_, err := r.DB.Exec(`INSERT INTO dbname (shortLink, longLink) VALUES ($1, $2)`,
		shortLink, longLink)

	return err
}

func(r *ItemRepostory) AddLink(longLink string) (string, error) {
	shortLink, err := r.searchShortLink(longLink)
	if err == nil {
		return shortLink, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	shortLink = GenerateShortLink()

	err = r.checkLink(longLink, shortLink)
	if err != nil {
		return "", err
	}

	err = r.addLinkInDB(shortLink, longLink)
	if err != nil {
		return "", err
	}

	return shortLink, nil
}
