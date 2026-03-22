package server

import (
	"log"
	"net"
)

func Start(port string) {
	//bind program to port and start accepting clients
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on tcp://0.0.0.0:6380")

	//handles multiple incoming connections
	for {
		conn, err := listener.Accept()
		log.Println("New connection", conn)
		if err != nil {
			log.Fatal(err)
		}
		
		//creates goroutine for each connection
		go startSession(conn)
	}
}