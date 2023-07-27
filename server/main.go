package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

var clients = make(map[net.Conn]struct{})

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer server.Close()

	fmt.Println("Server started, listening on port 8080")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Serving client: " + conn.RemoteAddr().String())

		clients[conn] = struct{}{}
		client := Client{conn}
		go handleClient(client)
	}
}

func handleClient(client Client) {
	defer client.conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := client.conn.Read(buf)
		if err != nil {
			delete(clients, client.conn)
			fmt.Println("Client disconnected")
			break
		}

		message := client.conn.RemoteAddr().String() + ": " + string(buf[:n])
		broadcastMessage(message, client.conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	for client := range clients {
		if client != sender {
			_, err := client.Write([]byte(message))
			if err != nil {
				fmt.Println("Error broadcasting message:", err)
			}
		}
	}
}

