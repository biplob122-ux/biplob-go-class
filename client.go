package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const SERVER_ADDRESS = "127.0.0.1:8080"

// receiveMessages continuously receives messages
// from the server.
func receiveMessages(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {

		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("\n[Server] Connection closed.")
			return
		}

		fmt.Print("\n" + message)
		fmt.Print("> ")
	}
}

func main() {

	fmt.Println("====================================")
	fmt.Println("          Go Chat Client")
	fmt.Println("====================================")
	fmt.Println()

	// Connect to the server.
	conn, err := net.Dial("tcp", SERVER_ADDRESS)

	if err != nil {
		fmt.Println("Unable to connect to server:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connected to server.")
	fmt.Println()

	// Create a reader for keyboard input.
	inputReader := bufio.NewReader(os.Stdin)

	// Ask the user for their name locally.
	fmt.Print("Enter your name: ")

	name, err := inputReader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}

	name = strings.TrimSpace(name)

	if name == "" {
		name = "Anonymous"
	}

	// Send the name to the server.
	_, err = fmt.Fprintln(conn, name)

	if err != nil {
		fmt.Println("Error sending name:", err)
		return
	}

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println("Connected to chat!")
	fmt.Println("Commands:")
	fmt.Println("  /users - Show online users")
	fmt.Println("  /help  - Show help")
	fmt.Println("  /quit  - Leave chat")
	fmt.Println("====================================")
	fmt.Println()

	// Start a goroutine to receive messages.
	go receiveMessages(conn)

	// Main loop for sending messages.
	for {

		fmt.Print("> ")

		message, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		message = strings.TrimSpace(message)

		if message == "" {
			continue
		}

		// Send message to server.
		_, err = fmt.Fprintln(conn, message)

		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		// Exit the client.
		if message == "/quit" {
			fmt.Println("Leaving chat...")
			break
		}
	}
}