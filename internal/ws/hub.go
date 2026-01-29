package ws

import (
	"sync"
)

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

func (h *Hub) handleMessage(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Send to specific user (private message)
	if msg.ReceiverID > 0 {
		if client, ok := h.clients[msg.ReceiverID]; ok {
			select {
			case client.send <- msg:
			default:
				// Client buffer full, skip
			}
		}
		return
	}

	// Broadcast to group members
	if msg.GroupID > 0 && len(msg.GroupMembers) > 0 {
		for _, memberID := range msg.GroupMembers {
			if memberID == msg.SenderID {
				continue // Don't send to sender
			}
			if client, ok := h.clients[memberID]; ok {
				select {
				case client.send <- msg:
				default:
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}

func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
