package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// Structs for handling WebSocket messages
type Message struct {
	Type          string   `json:"type"`
	Username      string   `json:"username"`
	Text          string   `json:"text,omitempty"`
	RoomID        string   `json:"room_id,omitempty"`
	ActiveMembers []string `json:"active_members,omitempty"`
}

// Struct for the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Global variables
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
	// Initialize the user table in the SQLite database
	createUserTable()

	// WebSocket and other routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections)

	// Add the login route
	http.HandleFunc("/login", handleLogin)

	go handleMessages()

	// Start the HTTP server
	log.Println("Server started on :7777")
	log.Fatal(http.ListenAndServe(":7777", nil))
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

// Helper function to create the SQLite users table if not exists
func createUserTable() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		push_data TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Users table created or already exists.")
}

// Helper function to validate the username (basic alphanumeric + underscore)
func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(username)
}

// Route for handling login and user creation
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON body for the login request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the username
	if !isValidUsername(req.Username) {
		http.Error(w, "Invalid username format", http.StatusBadRequest)
		return
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the user exists in the database
	var existingPassword, pushData string
	var userExists bool

	// Query to check if the user exists
	err = db.QueryRow("SELECT password, push_data FROM users WHERE username = ?", req.Username).Scan(&existingPassword, &pushData)
	if err == sql.ErrNoRows {
		// User does not exist, create a new user
		_, err := db.Exec("INSERT INTO users (username, password, push_data) VALUES (?, ?, ?)", req.Username, req.Password, "")
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		// Respond with success
		json.NewEncoder(w).Encode(LoginResponse{Success: true})
		return
	} else if err != nil {
		log.Fatal(err)
	}

	// Check if the password matches
	if existingPassword == req.Password {
		// Login success
		json.NewEncoder(w).Encode(LoginResponse{Success: true})
	} else {
		// Incorrect password
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
	}
}
