package command

import (
	"fmt"
	"miniredis/internal/pubsub"
)

func (cmd Command) publish() bool {
	if len(cmd.Args) != 3 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	count := pubsub.Publish(cmd.Args[1], cmd.Args[2])
	cmd.Conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))

	return true
}
