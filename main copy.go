package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type User struct {
	ID   int
	Name string
}

type Room struct {
	ID   int
	Name string
}

func initDB() {
	var err error
	// Open database connection
	db, err = sql.Open("sqlite3", "./chatapp.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS room_members (
		user_id INTEGER,
		room_id INTEGER,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, room_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);

	CREATE TABLE IF NOT EXISTS msgs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER,
		author_id INTEGER,
		message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id),
		FOREIGN KEY (author_id) REFERENCES users(id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

// Check if user has access to a specific room
func userCanJoinRoom(userID, roomID int) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM room_members WHERE user_id = ? AND room_id = ?`, userID, roomID).Scan(&count)
	if err != nil {
		log.Println("Error checking user access:", err)
		return false
	}
	return count > 0
}

// WebSocket handler
func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Read initial message containing user ID and room ID
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	// Assume the message is formatted as "userID-roomID"
	parts := strings.Split(string(msg), "-")
	if len(parts) != 2 {
		log.Println("Invalid message format")
		return
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println("Invalid userID format:", err)
		return
	}

	roomID, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Println("Invalid roomID format:", err)
		return
	}

	// Validate user and room
	if !userCanJoinRoom(userID, roomID) {
		conn.WriteMessage(websocket.TextMessage, []byte("You do not have access to this room."))
		return
	}

	// User can join the room, send a welcome message
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome to room %d", roomID)))

	// Handle incoming messages from the client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}
		log.Printf("Received message: %s", msg)

		// Send message back to client (this is just a sample echo)
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

// Route to join a room
func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Assume the userID and roomID are passed as URL parameters
	userIDStr := r.URL.Query().Get("user_id")
	roomIDStr := r.URL.Query().Get("room_id")

	if userIDStr == "" || roomIDStr == "" {
		http.Error(w, "User ID and Room ID are required", http.StatusBadRequest)
		return
	}

	// Convert userID and roomID to integers
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid Room ID", http.StatusBadRequest)
		return
	}

	// Check if room exists
	var roomCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM rooms WHERE id = ?`, roomID).Scan(&roomCount)
	if err != nil {
		http.Error(w, "Error checking room existence", http.StatusInternalServerError)
		return
	}
	if roomCount == 0 {
		http.Error(w, "Room does not exist", http.StatusNotFound)
		return
	}

	// Check if user can join the room
	if !userCanJoinRoom(userID, roomID) {
		http.Error(w, "You do not have access to this room", http.StatusForbidden)
		return
	}

	// Join room and add user to room_members table with joined_at timestamp
	_, err = db.Exec(`INSERT INTO room_members (user_id, room_id) VALUES (?, ?)`, userID, roomID)
	if err != nil {
		http.Error(w, "Error joining room", http.StatusInternalServerError)
		return
	}

	// Room join success
	w.Write([]byte(fmt.Sprintf("User %d successfully joined room %d", userID, roomID)))
}

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Setup routes
	http.HandleFunc("/ws", handleConnection)
	http.HandleFunc("/join_room", joinRoomHandler)

	// Start the HTTP server
	log.Println("Starting WebSocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
