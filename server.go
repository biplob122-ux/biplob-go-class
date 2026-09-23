package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	PORT        = ":8080"
	MAX_CLIENTS = 100
)

// Client represents one connected client.
type Client struct {
	conn net.Conn
	name string
}

// Server represents the chat server.
type Server struct {
	clients map[net.Conn]*Client
	mutex   sync.Mutex
}

// Creates a new Server.
func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]*Client),
	}
}

// Adds a client to the client list.
func (s *Server) addClient(client *Client) bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.clients) >= MAX_CLIENTS {
		return false
	}

	s.clients[client.conn] = client

	return true
}

// Removes a client from the client list.
func (s *Server) removeClient(conn net.Conn) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.clients, conn)
}

// Returns the number of connected clients.
func (s *Server) getClientCount() int {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.clients)
}

// Returns all online usernames.
func (s *Server) getUsers() string {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	var users []string

	for _, client := range s.clients {
		users = append(users, client.name)
	}

	return strings.Join(users, ", ")
}

// Sends a message to every client except the sender.
func (s *Server) broadcast(message string, sender net.Conn) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.clients {

		if conn != sender {
			fmt.Fprint(conn, message)
		}
	}
}

// Handles one client.
func (s *Server) handleClient(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)

name, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	name = strings.TrimSpace(name)

	if name == "" {
		name = "Anonymous"
	}

	client := &Client{
		conn: conn,
		name: name,
	}

	// Add client to shared client list.
	if !s.addClient(client) {

		fmt.Fprintln(
			conn,
			"[Server] Chat server is full.",
		)

		return
	}

	// Make sure client is removed when function ends.
	defer s.removeClient(conn)

	fmt.Printf(
		"[+] %s joined the chat. Total clients: %d\n",
		name,
		s.getClientCount(),
	)

	// Inform other clients.
	joinMessage := fmt.Sprintf(
		"[Server] %s joined the chat.\n",
		name,
	)

	s.broadcast(joinMessage, conn)

	// Welcome the new client.
	fmt.Fprintln(
		conn,
		"[Server] Welcome to the chat!",
	)

	fmt.Fprintln(
		conn,
		"[Server] Commands: /users, /help, /quit",
	)

	// Keep receiving messages.
	for {

		message, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		// Quit command.
		if message == "/quit" {
			break
		}

		// Show online users.
		if message == "/users" {

			users := s.getUsers()

			fmt.Fprintf(
				conn,
				"[Server] Online users: %s\n",
				users,
			)

			continue
		}

		// Help command.
		if message == "/help" {

			fmt.Fprintln(
				conn,
				"[Server] Available commands:",
			)

			fmt.Fprintln(
				conn,
				"  /users - Show online users",
			)

			fmt.Fprintln(
				conn,
				"  /help  - Show help",
			)

			fmt.Fprintln(
				conn,
				"  /quit  - Leave the chat",
			)

			continue
		}

		// Create chat message.
		chatMessage := fmt.Sprintf(
			"%s: %s\n",
			name,
			message,
		)

		// Display message on server.
		fmt.Print(chatMessage)

		// Send message to other clients.
		s.broadcast(chatMessage, conn)
	}

	// Inform other clients that this client left.
	leaveMessage := fmt.Sprintf(
		"[Server] %s left the chat.\n",
		name,
	)

	s.broadcast(leaveMessage, conn)

	fmt.Printf(
		"[-] %s left the chat. Total clients: %d\n",
		name,
		s.getClientCount()-1,
	)
}

// Starts the server.
func (s *Server) start() {

	listener, err := net.Listen(
		"tcp",
		PORT,
	)

	if err != nil {

		fmt.Println(
			"Error starting server:",
			err,
		)

		return
	}

	defer listener.Close()

	fmt.Println("====================================")
	fmt.Println("       Go TCP Chat Server")
	fmt.Println("====================================")
	fmt.Println("Server started on port 8080")
	fmt.Println("Waiting for clients...")
	fmt.Println()

	for {

		conn, err := listener.Accept()

		if err != nil {

			fmt.Println(
				"Error accepting connection:",
				err,
			)

			continue
		}

		fmt.Println(
			"[+] New client:",
			conn.RemoteAddr(),
		)

		// Create a goroutine for this client.
		go s.handleClient(conn)
	}
}

func main() {

	server := NewServer()

	server.start()
}