package item

import (
	"database/sql"
	"fmt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestSearchLongLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	longLink := "https://www.youtube.com/"
	shorLink := "as231_oskf"
	rows := sqlmock.NewRows([]string{"long_link"})
	expect := []*Item{&Item{ShortLink: shorLink, LongLink: longLink}}

	for _, item := range expect {
		rows = rows.AddRow(item.LongLink)
	}

	mock.
		ExpectQuery("SELECT long_link FROM dbname WHERE").
		WithArgs(shorLink).
		WillReturnRows(rows)

	repo := NewItemRepository(db)

	currLongLink, err := repo.SearchLongLink(shorLink)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if currLongLink != longLink {
		t.Errorf("results not match, want %v, have %v", longLink, currLongLink)
		return
	}
}

func TestSearchLongLinkError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT long_link FROM dbname WHERE").
		WithArgs("shorLink").
		WillReturnError(fmt.Errorf("db error"))

	repo := NewItemRepository(db)

	_, err = repo.SearchLongLink("shorLink")

	if err2 := mock.ExpectationsWereMet(); err2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestAddLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	longLink := "https://www.youtube.com/"
	shortLink := "asdf_29374"

	mock.
		ExpectQuery("SELECT short_link FROM dbname WHERE").
		WithArgs(longLink).WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery("SELECT long_link FROM dbname WHERE").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("INSERT INTO dbname").
		WithArgs(longLink, shortLink).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewItemRepository(db)

	GenerateShortLink = func() string {
		return shortLink
	}

	currShortLink, err := repo.AddLink(longLink)
	if err != nil || len(currShortLink) != 10  || currShortLink != shortLink{
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}


	// we have in db such row
	rows := sqlmock.NewRows([]string{"short_link"})
	expect := []*Item{&Item{ShortLink: shortLink, LongLink: longLink}}

	for _, item := range expect {
		rows = rows.AddRow(item.ShortLink)
	}

	mock.
		ExpectQuery("SELECT short_link FROM dbname WHERE").
		WithArgs(longLink).WillReturnRows(rows)

	repo = NewItemRepository(db)

	GenerateShortLink = func() string {
		return shortLink
	}

	currShortLink, err = repo.AddLink(longLink)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if currShortLink != shortLink {
		t.Errorf("results not match, want %v, have %v", shortLink, currShortLink)
		return
	}
}

func TestAddLinkError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	longLink := "https://www.youtube.com/"
	longLink2 := "https://github.com/"
	shortLink := "asdf_29374"

	// error in searchShortLink(not sql.ErrNoRows)
	mock.
		ExpectQuery("SELECT short_link FROM dbname WHERE").
		WithArgs(longLink).WillReturnError(fmt.Errorf("error db"))

	repo := NewItemRepository(db)

	_, err = repo.AddLink(longLink)

	if err2 := mock.ExpectationsWereMet(); err2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// checkLink error
	mock.
		ExpectQuery("SELECT short_link FROM dbname WHERE").
		WithArgs(longLink).WillReturnError(sql.ErrNoRows)

	rows := sqlmock.NewRows([]string{"short_link"})
	expect := []*Item{&Item{ShortLink: shortLink, LongLink: longLink2}}

	for _, item := range expect {
		rows = rows.AddRow(item.LongLink)
	}

	mock.
		ExpectQuery("SELECT long_link FROM dbname WHERE").
		WithArgs(shortLink).WillReturnRows(rows)

	repo = NewItemRepository(db)

	GenerateShortLink = func() string {
		return shortLink
	}
	_, err = repo.AddLink(longLink)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err2 := mock.ExpectationsWereMet(); err2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}



	// insert error
	mock.
		ExpectQuery("SELECT short_link FROM dbname WHERE").
		WithArgs(longLink).WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery("SELECT long_link FROM dbname WHERE").
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("INSERT INTO dbname").
		WithArgs(longLink, shortLink).WillReturnError(fmt.Errorf("db error"))

	repo = NewItemRepository(db)

	GenerateShortLink = func() string {
		return shortLink
	}
	_, err = repo.AddLink(longLink)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err2 := mock.ExpectationsWereMet(); err2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

