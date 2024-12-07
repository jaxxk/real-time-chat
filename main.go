package main

import (
	"log"
	"net/http"

	server "github.com.jaxxk.real-time-chat/cmd/chat-server"
	"github.com.jaxxk.real-time-chat/logger"
)

var addr = "localhost:8080"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	logger.Init("server.log")
	chatServer := server.NewChatServer("")
	logger.Info.Println("Starting Chat Server...")

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(chatServer, w, r)
	})

	go chatServer.RunChatServer()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error.Fatalf("ListenAndServe failed: %v", err)
	}
}
