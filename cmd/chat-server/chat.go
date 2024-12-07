package server

import (
	"sync"

	"github.com.jaxxk.real-time-chat/logger"
)

type ChatServer struct {
	clients    sync.Map // Thread-safe map for clients
	messenger  Messenger
	register   chan *Client
	unregister chan *Client
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		messenger:  *newMessenger(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (c *ChatServer) RunChatServer() {
	for {
		select {
		case client := <-c.register:
			// Register a new client
			logger.Info.Printf("Registering client: %v", client.name)
			c.clients.Store(client, true)

		case client := <-c.unregister:
			// Unregister an existing client
			if _, ok := c.clients.Load(client); ok {
				logger.Info.Printf("Unregistering client: %v", client.name)
				c.clients.Delete(client)
				close(client.send) // Close the client's send channel
			}

		case message := <-c.messenger.broadcast:
			// Broadcast the message to all clients
			c.clients.Range(func(key, value interface{}) bool {
				client := key.(*Client) // Cast to *Client
				select {
				case client.send <- message:
				default:
					// If the client is unresponsive, remove it
					logger.Info.Printf("Removing unresponsive client: %v", client.name)
					c.clients.Delete(client)
					close(client.send)
				}
				return true // Continue iteration
			})
		}
	}
}
