package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com.jaxxk.real-time-chat/logger"
	"github.com/gorilla/websocket"
)

// Client represents a single chat client.
type Client struct {
	chat *ChatServer
	conn *websocket.Conn
	name string
	send chan []byte
}

const (
	writeWait  = 10 * time.Second    // Time allowed to write a message to the peer
	pongWait   = 60 * time.Second    // Time allowed to read the next pong message from the peer
	pingPeriod = (pongWait * 9) / 10 // Send pings to peer periodically
	maxMsgSize = 512                 // Maximum message size allowed from peer
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// readPump listens for messages from the WebSocket connection and broadcasts them.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.chat.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			logger.Info.Printf("Stopping readPump for client %s: %v", c.name, ctx.Err())
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error.Printf("Unexpected close error for client %s: %v", c.name, err)
				}
				return
			}

			message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
			formattedMessage := fmt.Sprintf("%s: %s", c.name, string(message))
			c.chat.messenger.broadcast <- []byte(formattedMessage)
		}
	}
}

// writePump sends messages from the server to the WebSocket client.
func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Info.Printf("Stopping writePump for client %s: %v", c.name, ctx.Err())
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel closed
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			if err := c.writeMessageWithQueue(message); err != nil {
				return
			}

		case <-ticker.C:
			// Send a ping message to check connection health
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// writeMessageWithQueue writes a single message to the WebSocket, including queued messages.
func (c *Client) writeMessageWithQueue(message []byte) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	defer w.Close()

	// Write the initial message
	if _, err := w.Write(message); err != nil {
		return err
	}

	// Write any additional queued messages
	n := len(c.send)
	for i := 0; i < n; i++ {
		w.Write(newline)
		w.Write(<-c.send)
	}

	return nil
}

// ServeWs handles WebSocket requests from the peer.
func ServeWs(chatServer *ChatServer, w http.ResponseWriter, r *http.Request, ctx context.Context) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		chat: chatServer,
		conn: conn,
		name: name,
		send: make(chan []byte, 256),
	}

	// Register the client and start its read/write pumps
	client.chat.register <- client

	go client.readPump(ctx)
	go client.writePump(ctx)

}
