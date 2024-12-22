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
	RoomID   string `json:"room_id,omitempty"`
}

var (
	clients   = make(map[*websocket.Conn]string)          // Track client connections and their rooms
	rooms     = make(map[string]map[*websocket.Conn]bool) // Rooms with connected clients
	broadcast = make(chan Message)                        // Broadcast channel
	upgrader  = websocket.Upgrader{                       // WebSocket upgrader
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
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Extract room ID from the URL query parameter
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		roomID = "default" // If no room ID is provided, assign to default room
	}

	// Assign client to room
	clients[conn] = roomID
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID][conn] = true

	// Listen for incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(rooms[roomID], conn)
			delete(clients, conn)
			break
		}
		msg.RoomID = roomID // Set room ID for the message
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		// Broadcast message only to clients in the same room
		for client := range rooms[msg.RoomID] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(rooms[msg.RoomID], client)
				delete(clients, client)
			}
		}
	}
}
