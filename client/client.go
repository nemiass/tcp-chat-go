package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	protocol = "tcp"
	port     = ":9999"
)

var username string

func main() {

	serverConnection, err := net.Dial(protocol, port)
	defer func() {
		serverConnection.Close()
	}()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Develop by Nemias")
	fmt.Print("username: ")
	fmt.Scan(&username)

	err = gob.NewEncoder(serverConnection).Encode(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	go receive_message(serverConnection)
	write_message(serverConnection)
}

func receive_message(serverConnection net.Conn) {
	msgFromServer := make([]byte, 100)

	for {
		bytesReaded, err := serverConnection.Read(msgFromServer)
		if err != nil {
			log.Println("Server desconceted!!")
			break
		}

		if strings.Contains(string(msgFromServer[:bytesReaded]), "ya existe") {
			fmt.Println(string(msgFromServer[:bytesReaded]))
			continue
		}
		//fmt.Println(string(msgFromServer[:bytesReaded]) + "\r")
		fmt.Printf("\r%s\n", string(msgFromServer[:bytesReaded]))
	}
}

func write_message(serverConnection net.Conn) {
	//reader := bufio.NewReader(os.Stdin)
	sc := bufio.NewScanner(os.Stdin)
	for {
		//message, _ := reader.ReadString('\n')
		sc.Scan()
		message := sc.Text()
		if len(message) == 1 {
			continue
		}
		_, err := serverConnection.Write([]byte(message))
		if err != nil {
			log.Println("Server not Responding!!")
			break
		}
	}
}
