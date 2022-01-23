package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"link_app/pkg/item"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestItemHandler_GetLongLink(t *testing.T) {
	// мы передаём t сюда, это надо, чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemRepositoryInterface(ctrl)
	service := &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	longLink := "https://www.youtube.com/"
	shortLink := "shortLink1"

	st.EXPECT().SearchLongLink(shortLink).Return(longLink, nil)

	req := httptest.NewRequest("GET", "/link/shortLink1", nil)
	vars := map[string]string{
		"SHORT_LINK": "shortLink1",
	}
	req = mux.SetURLVars(req, vars)
	w := httptest.NewRecorder()

	service.GetLongLink(w, req)

	resp := w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	//nolint:errcheck
	body, _ := ioutil.ReadAll(resp.Body)

	us := item.Item{}

	err := json.Unmarshal(body, &us)
	if err != nil || longLink != us.LongLink {
		t.Errorf("results not match, want %v, have %v", longLink, us.LongLink)
		return
	}

	// error
	st = NewMockItemRepositoryInterface(ctrl)
	service = &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}
	st.EXPECT().SearchLongLink(shortLink).Return("", fmt.Errorf("error"))

	req = mux.SetURLVars(req, vars)
	w = httptest.NewRecorder()

	service.GetLongLink(w, req)
	resp = w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestItemHandler_CreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemRepositoryInterface(ctrl)
	service := &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	longLink := "https://www.youtube.com/"
	shortLink := "shortLink1"
	resultItem := &item.Item{
		LongLink:  longLink,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))
	item.GenerateShortLink = func () string{return shortLink}
	st.EXPECT().AddLink(longLink).Return(shortLink, nil)

	req := httptest.NewRequest("POST", "/add", bodyReader)
	w := httptest.NewRecorder()

	//s := service.CreateShortLink1(generate)
	//s(w, req)
	service.CreateShortLink(w, req)

	resp := w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	//nolint:errcheck
	body, _ := ioutil.ReadAll(resp.Body)

	us := item.Item{}

	err = json.Unmarshal(body, &us)
	if err != nil || shortLink != us.ShortLink {
		t.Errorf("results not match, want %v, have %v", longLink, us.LongLink)
		return
	}

}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("test error")
}


func TestItemHandler_CreateShortLinkError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := NewMockItemRepositoryInterface(ctrl)
	service := &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	req := httptest.NewRequest("POST", "/add", errReader(0))
	w := httptest.NewRecorder()


	service.CreateShortLink(w, req)

	resp := w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	//nolint:errcheck
	if resp.StatusCode != 400 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	//marshal error
	st = NewMockItemRepositoryInterface(ctrl)
	service = &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	bodyReader := strings.NewReader(`short`)
	req = httptest.NewRequest("POST", "/api/login", bodyReader)
	w = httptest.NewRecorder()

	service.CreateShortLink(w, req)

	resp = w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	//nolint:errcheck
	if resp.StatusCode != 400 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	// create error
	st = NewMockItemRepositoryInterface(ctrl)
	service = &ItemHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	longLink := "https://www.youtube.com/"
	shortLink := "shortLink1"
	resultItem := &item.Item{
		LongLink:  longLink,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader = strings.NewReader(string(b))
	item.GenerateShortLink = func () string{return shortLink}
	st.EXPECT().AddLink(longLink).Return("", fmt.Errorf("error"))

	req = httptest.NewRequest("POST", "/api/login", bodyReader)
	w = httptest.NewRecorder()
	service.CreateShortLink(w, req)

	resp = w.Result()
	//nolint:errcheck
	defer resp.Body.Close()

	//nolint:errcheck
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}
