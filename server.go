package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
)

var users map[string]net.Conn

var (
	protocol = "tcp"
	port     = ":9999"
)

var mu sync.Mutex

func main() {

	users = make(map[string]net.Conn)

	listener, err := net.Listen(protocol, port)

	defer func() {
		listener.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	var username string
	fmt.Println("Server UP!! in port:" + port)
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = gob.NewDecoder(client).Decode(&username)
		if err != nil {
			fmt.Println(err)
			client.Close()
			continue
		}

		if _, ok := users[username]; ok {
			err = gob.NewEncoder(client).Encode(username + " ya existe!")
			if err != nil {
				fmt.Println(err)
			}
			client.Close()
			continue
		}

		users[username] = client
		log.Println("User conected: " + username)
		msgAll := "<" + username + "> joined in the server!"
		broadcasting(msgAll)
		go start_session(client, username)
	}
}

func start_session(client net.Conn, username string) {

	defer func() {
		client.Close()
		mu.Lock()
		delete(users, username)
		mu.Unlock()
	}()

	msgFromClient := make([]byte, 100)

	for {
		bytesRead, err := client.Read(msgFromClient)
		if err != nil {
			log.Println("User desconected: " + username)
			broadcasting("<" + username + "> left to the server!!")
			break
		}
		broadcasting("<" + username + ">: " + string(msgFromClient[:bytesRead]))
	}
}

func broadcasting(msg string) {
	mu.Lock()
	for _, client := range users {

		_, err := client.Write([]byte(msg))
		if err != nil {
			continue
		}
	}
	mu.Unlock()
}
