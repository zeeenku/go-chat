package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type          string   `json:"type"`
	Username      string   `json:"username"`
	Text          string   `json:"text,omitempty"`
	RoomID        string   `json:"room_id,omitempty"`
	ActiveMembers []string `json:"active_members,omitempty"` // List of active members
}

var (
	clients   = make(map[*websocket.Conn]string)          // Track client connections and their usernames
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

	// Extract room ID and username from URL query parameters
	roomID := r.URL.Query().Get("room_id")
	username := r.URL.Query().Get("username") // Ensure that the client sends the username
	if roomID == "" {
		roomID = "default" // If no room ID is provided, assign to default room
	}
	if username == "" {
		username = "Anonymous" // Default username if not provided
	}

	// Assign client to room and store their username
	clients[conn] = username
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID][conn] = true

	// Notify all users about active members
	updateActiveMembers(roomID)

	// Listen for incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(rooms[roomID], conn)
			delete(clients, conn)
			updateActiveMembers(roomID) // Update active members after a user leaves
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

// Update active members list in the room
func updateActiveMembers(roomID string) {
	activeMembers := []string{}
	for client := range rooms[roomID] {
		activeMembers = append(activeMembers, clients[client]) // Use username instead of room ID
	}

	// Broadcast the active members list to all clients in the room
	for client := range rooms[roomID] {
		err := client.WriteJSON(Message{
			Type:          "active-members",
			RoomID:        roomID,
			ActiveMembers: activeMembers,
		})
		if err != nil {
			log.Printf("Error sending active members update: %v", err)
		}
	}
}
