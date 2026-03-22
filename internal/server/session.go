package server

import (
	"log"
	"miniredis/internal/parser"
	"net"
)

// startSession handles the client's session. Parses and executes commands and writes
// responses back to the client.
func startSession(conn net.Conn) {

	//normal connection cleanup
	defer func() {
		log.Println("Closing connection", conn)
		conn.Close()
	}()

	//helps server not crash due to unexpected panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovering from error", err)
		}
	}()

	//initialize parser
	p := parser.NewParser(conn)
	for {
		//continuously reads from client
		cmd, err := p.Command()

		//returns "-ERR" incase of parsing failure
		if err != nil {
			log.Println("Error", err)
			conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
			break
		}

		//if cmd (get, set etc) return false, session ends
		if !cmd.Handle() {
			break
		}
	}
}
