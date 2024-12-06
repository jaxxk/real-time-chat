package server

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
