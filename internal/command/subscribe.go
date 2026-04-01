package command

import (
	"miniredis/internal/pubsub"
)

func (cmd Command) subscribe() bool {
	if len(cmd.Args) < 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	pubsub.Subscribe(cmd.Args[1], cmd.Conn)

	return false
}
