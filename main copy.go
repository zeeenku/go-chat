package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB
var clients = make(map[*websocket.Conn]string)        // Map of connected clients
var rooms = make(map[string]map[*websocket.Conn]bool) // Map of rooms and their clients
var broadcast = make(chan Message)                    // Channel for broadcasting messages
var mu sync.Mutex

type Message struct {
	RoomID   string `json:"room_id"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

func initDB() {
	var err error
	// Open the database connection
	db, err = sql.Open("sqlite3", "./chatapp.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the tables if they don't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS room_members (
		user_id INTEGER,
		room_id INTEGER,
		PRIMARY KEY (user_id, room_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER,
		username TEXT,
		text TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func addUser(username string) (int, error) {
	var userID int
	err := db.QueryRow(`INSERT INTO users (username) VALUES (?) RETURNING id`, username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func addRoom(name string) (int, error) {
	var roomID int
	err := db.QueryRow(`INSERT INTO rooms (name) VALUES (?) RETURNING id`, name).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func addUserToRoom(userID, roomID int) error {
	_, err := db.Exec(`INSERT INTO room_members (user_id, room_id) VALUES (?, ?)`, userID, roomID)
	return err
}

func saveMessage(roomID int, username, text string) error {
	_, err := db.Exec(`INSERT INTO messages (room_id, username, text) VALUES (?, ?, ?)`, roomID, username, text)
	return err
}

func updateActiveMembers(roomID string) {
	mu.Lock()
	defer mu.Unlock()
	// Send a list of active members in the room to all clients
	var users []string
	rows, err := db.Query(`SELECT username FROM users WHERE id IN (SELECT user_id FROM room_members WHERE room_id = ?)`, roomID)
	if err != nil {
		log.Printf("Error retrieving active members: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			log.Printf("Error scanning active member: %v", err)
			continue
		}
		users = append(users, username)
	}

	// Broadcast active members to all clients in the room
	for conn := range rooms[roomID] {
		err := conn.WriteJSON(Message{
			RoomID:   roomID,
			Username: "System",
			Text:     fmt.Sprintf("Active members: %v", users),
		})
		if err != nil {
			log.Printf("Error sending active members to %v: %v", conn.RemoteAddr(), err)
		}
	}
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
	username := r.URL.Query().Get("username")
	if roomID == "" {
		roomID = "default" // Default room
	}
	if username == "" {
		username = "Anonymous" // Default username
	}

	// Check if the user exists, otherwise add them
	var userID int
	err = db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		// User doesn't exist, so add them
		userID, err = addUser(username)
		if err != nil {
			log.Printf("Error adding user: %v", err)
			return
		}
	}

	// Add user to the room
	roomIDInt, err := strconv.Atoi(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		return
	}
	err = addUserToRoom(userID, roomIDInt)
	if err != nil {
		log.Printf("Error adding user to room: %v", err)
		return
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

		// Save the message to the database
		err = saveMessage(roomIDInt, username, msg.Text)
		if err != nil {
			log.Printf("Error saving message to database: %v", err)
		}

		// Broadcast message to all clients in the same room
		broadcast <- msg
	}
}

func handleMessages() {
	for msg := range broadcast {
		mu.Lock()
		for conn := range rooms[msg.RoomID] {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to %v: %v", conn.RemoteAddr(), err)
				conn.Close()
				delete(rooms[msg.RoomID], conn)
			}
		}
		mu.Unlock()
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// Initialize database
	initDB()

	// Set up WebSocket route
	http.HandleFunc("/ws", handleConnections)

	// Start message handler
	go handleMessages()

	// Start the server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
