package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"link_app/pkg/item"
	"net/http"
)

//go:generate mockgen -source=item.go -destination=item_mock.go -package=handler ItemRepositoryInterface
type ItemRepositoryInterface interface {
	SearchLongLink(shortLink string) (string, error)
	AddLink(longLink string) (string, error)
}


type ItemHandler struct {
	Logger   *zap.SugaredLogger
	ItemRepo ItemRepositoryInterface
}

func SendDataHandler(w http.ResponseWriter, r *http.Request, data interface{}, logger *zap.SugaredLogger) {
	dataJSON, _ := json.Marshal(data)
	w.Header().Set("Content-type", "application/json")
	w.Write(dataJSON)
}


func (h *ItemHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendErrorHandler(w, r, err,
			http.StatusBadRequest, h.Logger)
		return
	}

	currItem := &item.Item{}
	err = json.Unmarshal(b, &currItem)
	if err != nil {
		SendErrorHandler(w, r, err,
			http.StatusBadRequest, h.Logger)
		return
	}

	shortLink, err := h.ItemRepo.AddLink(currItem.LongLink)
	if err != nil {
		SendErrorHandler(w, r, err,
			http.StatusInternalServerError, h.Logger)
		return
	}

	tmpLink := item.Item{ShortLink: shortLink}

	SendDataHandler(w, r, tmpLink, h.Logger)
}


func (h *ItemHandler) GetLongLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["SHORT_LINK"]

	longLink, err := h.ItemRepo.SearchLongLink(shortLink)
	if err != nil {
		SendErrorHandler(w, r, err,
			http.StatusInternalServerError, h.Logger)
		return
	}

	currLink := item.Item{LongLink: longLink}

	SendDataHandler(w, r, currLink, h.Logger)
}
