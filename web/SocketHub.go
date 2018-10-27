package web

import (
	"encoding/json"
	"time"
)

// SocketHub maintains the set of active clients and broadcasts messages to them
type SocketHub struct {
	// Registered clients.
	clients map[*SocketClient]bool

	// Inbound messages from the clients.
	broadcastChannel chan []byte

	// Register requests from the clients.
	registerChannel chan *SocketClient

	// Unregister requests from clients.
	unregisterChannel chan *SocketClient
}

func newSocketHub() *SocketHub {
	return &SocketHub{
		broadcastChannel:  make(chan []byte),
		registerChannel:   make(chan *SocketClient),
		unregisterChannel: make(chan *SocketClient),
		clients:           make(map[*SocketClient]bool),
	}
}

const (
	messageInterval = 2 * time.Second
)

// SocketMessage encapsulates a message to the client
type SocketMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (hub *SocketHub) broadcast(typ string, data string) {
	go func() {
		message := &SocketMessage{Type: typ, Data: data}
		jsonBytes, _ := json.Marshal(message)
		hub.broadcastChannel <- jsonBytes
	}()
}

func (hub *SocketHub) run() {
	for {
		select {
		case client := <-hub.registerChannel:
			hub.clients[client] = true

		case client := <-hub.unregisterChannel:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}

		case message := <-hub.broadcastChannel:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
