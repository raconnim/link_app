package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string `json:"message,omitempty"`
}

func SendErrorHandler(w http.ResponseWriter, r *http.Request, err error, status int, logger *zap.SugaredLogger) {
	logger.Errorw("SendErrorHandler",
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
		"url", r.URL.Path,
		"time", time.Now(),
		"err", err.Error())
	resp := ErrorResponse{Error: err.Error()}

	dataJSON, _ := json.Marshal(resp)
	http.Error(w, string(dataJSON), status)
}
