package utils

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"sync"
)

type Client struct {
	ID      int
	Conn    *websocket.Conn
	Message chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
	mu         sync.Mutex
	log        zerolog.Logger
}

func NewHub(logger zerolog.Logger) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		log:        logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
		case client := <-h.UnRegister:
			h.mu.Lock()
			delete(h.Clients, client)
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				select {
				case client.Message <- message:
				default:
					close(client.Message)
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
