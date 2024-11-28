package handler

import (
	"client/request"
	"encoding/json"
	"net/http"
)

type Handler struct {
	httpRequest *request.HttpRequest
}

func NewHandler(httpRequest *request.HttpRequest) *Handler {
	return &Handler{httpRequest}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {

	res, err := h.httpRequest.PingToServeApp()
	data := map[string]string{
		"message": "Message",
		"data":    string(res),
	}

	byteData, _ := json.Marshal(data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(byteData)
}
