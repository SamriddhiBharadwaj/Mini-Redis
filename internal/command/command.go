package command

import (
	"log"
	"net"
	"strings"
)

// Command implements the behavior of the commands.
type Command struct {
	Args []string
	Conn net.Conn
}

// handle Executes the command and writes the response. Returns false when the connection should be closed.
func (cmd Command) Handle() bool {
	switch strings.ToUpper(cmd.Args[0]) {
	case "GET":
		return cmd.get()
	case "SET":
		return cmd.set()
	case "DEL":
		return cmd.del()
	case "QUIT":
		return cmd.quit()
	case "TTL":
		return cmd.ttl()
	case "INCR":
		return cmd.incr()
	case "DECR":
		return cmd.decr()
	default:
		log.Println("Command not supported", cmd.Args[0])
		cmd.Conn.Write([]uint8("-ERR unknown command '" + cmd.Args[0] + "'\r\n"))
	}
	return true
}
