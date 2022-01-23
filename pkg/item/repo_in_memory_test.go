package item

import (
	"testing"
)

type TestCase struct {
	LongLink     string
	ShortLink string
	IsError bool
}

func TestItemRepositoryInMemory_AddLink(t *testing.T) {
	cases := []TestCase{
		{ LongLink:  "https://www.youtube.com/", ShortLink: "shortLink1", IsError:   false},
		{ LongLink:  "https://github.com/", ShortLink: "shortLink2", IsError:   false},
		{ LongLink:  "https://twitter.com/", ShortLink: "shortLink1", IsError:   true},
	}

	repo := NewItemRepositoryInMemory()
	repo.data[cases[0].ShortLink] = cases[0].LongLink
	for caseNum, item := range cases {
		GenerateShortLink = func()string{return item.ShortLink}
		shortLink, err := repo.AddLink(item.LongLink)
		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		if shortLink != item.ShortLink && err == nil {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, shortLink, item.ShortLink)
		}

	}
}
func TestItemRepositoryInMemory_SearchLongLink(t *testing.T) {

	cases := []TestCase{
		{ LongLink:  "https://www.youtube.com/", ShortLink: "shortLink1", IsError:   false},
		{ LongLink:  "https://github.com/", ShortLink: "shortLink2", IsError:   false},
		{ LongLink:  "https://twitter.com/", ShortLink: "", IsError:   true},
	}

	repo := NewItemRepositoryInMemory()
	repo.data[cases[0].ShortLink] = cases[0].LongLink
	repo.data[cases[1].ShortLink] = cases[1].LongLink

	for caseNum, item := range cases {
		longLink, err := repo.SearchLongLink(item.ShortLink)
		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		if longLink != item.LongLink && err == nil {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, longLink, item.LongLink)
		}

	}
}