package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type MSG struct {
	From    string
	To      string
	Message string
}

func main() {
	myName := "client1"

	conn, err := net.Dial("tcp", ":4545")
	if err != nil {
		fmt.Println("Connection lost: ", err)
		return
	}
	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	err = encoder.Encode(MSG{From: myName})
	if err != nil {
		fmt.Println("Failed to send name: ", err)
		return
	}

	go func() {
		for {
			var msg MSG
			if err := decoder.Decode(&msg); err != nil {
				fmt.Println("Disconnected or decode error: ", err)
				os.Exit(0)
			}
			fmt.Printf("\n%s> %s\n", msg.From, msg.Message)
		}
	}()
	for {
		var to string
		var message string
		to = "client"
		_, err = fmt.Scanln(&message)
		if err != nil {
			fmt.Println("Scan error: ", err)
		}
		msg := MSG{
			From:    myName,
			To:      to,
			Message: message,
		}
		if err := encoder.Encode(msg); err != nil {
			fmt.Println("Send failed: ", err)
			return
		}
	}
}
