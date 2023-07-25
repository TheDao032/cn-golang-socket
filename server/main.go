// package main
//
// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// )
//
// const (
// 	SERVER_HOST = "localhost"
// 	SERVER_PORT = "9988"
// 	SERVER_TYPE = "tcp"
// )
//
// func main() {
// 	fmt.Println("Server Running...")
// 	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
// 	if err != nil {
// 		fmt.Println("Error listening:", err.Error())
// 		os.Exit(1)
// 	}
// 	defer server.Close()
// 	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
// 	fmt.Println("Waiting for client...")
// 	for {
// 		connection, err := server.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting: ", err.Error())
// 			os.Exit(1)
// 		}
// 		fmt.Println("client connected")
// 		go processRead(connection)
// 		go processWrite(connection)
// 	}
// }
//
// func processRead(c net.Conn) {
// 	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
// 	for {
// 		buffer := make([]byte, 1024)
// 		mLen, err := c.Read(buffer)
// 		if err != nil {
// 			fmt.Println("Error reading:", err.Error())
// 			break
// 		}
// 		fmt.Println("Client: ", string(buffer[:mLen]))
// 		// netData, err := bufio.NewReader(c).ReadString('\n')
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return
// 		// }
//
// 		temp := strings.TrimSpace(string(buffer[:mLen]))
// 		if temp == "STOP" {
// 			break
// 		}
// 	}
// 	c.Close()
//
// 	// buffer := make([]byte, 1024)
// 	// mLen, err := connection.Read(buffer)
// 	// if err != nil {
// 	// 	fmt.Println("Error reading:", err.Error())
// 	// }
// 	// fmt.Println("Received: ", string(buffer[:mLen]))
// 	// _, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
// 	// connection.Close()
// }
//
// func processWrite(c net.Conn) {
// 	// buffer := make([]byte, 1024)
// 	// mLen, err := connection.Read(buffer)
// 	// if err != nil {
// 	// 	fmt.Println("Error reading:", err.Error())
// 	// }
// 	// fmt.Println("Received: ", string(buffer[:mLen]))
// 	// _, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
// 	// connection.Close()
// 	// fmt.Printf("Serving %s\n", c.RemoteAddr().String())
// 	for {
// 		fmt.Print("input text: ")
// 		netData, err := bufio.NewReader(c).ReadString('\n')
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
//
// 		if err != nil {
// 			fmt.Println("Error reading:", err.Error())
// 			continue
// 		}
//c
// 		temp := strings.TrimSpace(string(netData))
// 		if temp == "STOP" {
// 			break
// 		}
//
// 		fmt.Printf("Server: %s\n", temp)
// 		_, err = c.Write([]byte(temp))
// 	}
// 	c.Close()
// }

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

