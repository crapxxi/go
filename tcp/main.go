package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

type MSG struct {
	From    string
	To      string
	Message string
}

var (
	clients = make(map[string]chan MSG)
	mu      = sync.Mutex{}
)

func main() {
	gob.Register(MSG{})

	listener, err := net.Listen("tcp", ":4545")
	if err != nil {
		fmt.Println("Listen error: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error: ", err)
		}
		go handleClient(conn)
	}
}
func handleClient(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var fm MSG
	if err := decoder.Decode(&fm); err != nil {
		fmt.Println("Failed to read first message: ", err)
		return
	}
	clientName := fm.From
	clientChan := make(chan MSG, 10)
	mu.Lock()
	clients[clientName] = clientChan
	mu.Unlock()

	fmt.Println("Client registered: ", clientName)

	go func() {
		for msg := range clientChan {
			if err := encoder.Encode(msg); err != nil {
				fmt.Println("Write error: ", err)
				return
			}
		}
	}()
	for {
		var msg MSG
		if err := decoder.Decode(&msg); err != nil {
			fmt.Println("Decode err: ", err)
			break
		}

		mu.Lock()
		targetChan, ok := clients[msg.To]
		mu.Unlock()
		if ok {
			targetChan <- msg
		} else {
			fmt.Println("User not found: ", msg.To)
		}
	}

	mu.Lock()
	delete(clients, clientName)
	close(clientChan)
	mu.Unlock()

	fmt.Println("Client disconnected: ", clientName)
}
