package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Text     string `json:"text,omitempty"`
}

var (
	clients   = make(map[*websocket.Conn]bool) // Connected clients
	broadcast = make(chan Message)             // Broadcast channel
	upgrader  = websocket.Upgrader{            // WebSocket upgrader
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(clients, conn)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
