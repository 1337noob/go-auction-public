package application

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	EventType string `json:"event_type"`
	EventData any    `json:"event_data"`
	UserID    string `json:"user_id"`
	Broadcast bool   `json:"broadcast"`
}

type WsMessage struct {
	EventType string `json:"event_type"`
	EventData any    `json:"event_data"`
}

type Connection struct {
	UserID string
	Conn   *websocket.Conn
}

type WsHub struct {
	hub map[string]map[string]Connection
	mu  sync.Mutex
}

func NewWsHub() *WsHub {
	return &WsHub{
		hub: make(map[string]map[string]Connection),
	}
}

func (h *WsHub) AddConnection(auctionID string, userID string, conn *websocket.Conn) {
	log.Println("Adding connection to hub")
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.hub[auctionID]; !ok {
		h.hub[auctionID] = make(map[string]Connection)
	}
	h.hub[auctionID][userID] = Connection{
		UserID: userID,
		Conn:   conn,
	}
	log.Println("Added connection to hub")
}

func (h *WsHub) RemoveConnection(auctionID string, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.hub[auctionID], userID)
}

func (h *WsHub) getConnections(auctionID string) []Connection {
	connections := make([]Connection, 0, len(h.hub[auctionID]))
	for _, conn := range h.hub[auctionID] {
		connections = append(connections, conn)
	}

	return connections
}

func (h *WsHub) SendMessage(auctionID string, message *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections := h.getConnections(auctionID)

	for _, conn := range connections {
		if !message.Broadcast && message.UserID != conn.UserID {
			continue
		}

		//log.Println("Sending message to conn: ", conn.Conn.RemoteAddr().String())
		wsMessage := &WsMessage{
			EventType: message.EventType,
			EventData: message.EventData,
		}
		jsonMessage, _ := json.Marshal(wsMessage)
		err := conn.Conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Error writing message: %v\n", err)
			continue
		}
		log.Printf("Message sent: %v\n", message)
	}
}
