package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func handleIncomingMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Print(message)
	}
}

func main() {
	args := os.Args
	conn, err := net.Dial("tcp", args[1])
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()


	fmt.Println("connecting to server successfully")
	go handleIncomingMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading input:", scanner.Err())
	}
}
