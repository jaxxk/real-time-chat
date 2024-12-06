package server

import "github.com/gorilla/websocket"

type Client struct {
	chat *ChatServer
	conn *websocket.Conn
}
