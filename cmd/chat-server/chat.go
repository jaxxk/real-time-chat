package server

import "github.com.jaxxk.real-time-chat/logger"

type ChatServer struct {
	clients map[*Client]bool

	messenger Messenger

	register chan *Client

	unregister chan *Client
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:    make(map[*Client]bool),
		messenger:  *newMessenger(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (c ChatServer) RunChatServer() {
	for {
		select {

		case client := <-c.register:
			logger.Info.Printf("Registering client: %v", client.name)
			c.clients[client] = true

		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				logger.Info.Printf("Unregistering client: %v", client.name)
				delete(c.clients, client)
			}

		case message := <-c.messenger.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					// Remove clients that can't receive messages
					close(client.send)
					delete(c.clients, client)
				}
			}
		}

	}
}
