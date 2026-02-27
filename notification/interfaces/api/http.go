package api

import (
	"encoding/json"
	"log"
	"main/notification/module"
	"net/http"

	"github.com/gorilla/websocket"
)

type HttpHandler struct {
	module *module.NotificationModule
}

func NewHttpHandler(module *module.NotificationModule) *HttpHandler {
	return &HttpHandler{
		module: module,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *HttpHandler) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (h *HttpHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	auctionID := r.URL.Query().Get("auction_id")
	if auctionID == "" {
		h.writeError(w, http.StatusBadRequest, "auction_id is required")
		return
	}
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade websocket connection: %s", err)
		h.writeError(w, http.StatusInternalServerError, "failed to upgrade websocket connection: "+err.Error())
		return
	}

	h.module.AddConnection(auctionID, userID, conn)

	log.Println("Added connection to auction hub", auctionID)
}
