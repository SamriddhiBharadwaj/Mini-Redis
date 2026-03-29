package command

import (
	"fmt"
	"miniredis/internal/cache"
)

func (cmd Command) ttl() bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for 'TTL' command\r\n"))
		return true
	}

	ttl := cache.TTL(cmd.Args[1])

	cmd.Conn.Write([]uint8(fmt.Sprintf(":%d\r\n", int(ttl.Seconds()))))
	return true
}
