package command

import (
	"fmt"
	"log"
	"miniredis/internal/cache"
	"strconv"
	"strings"
)

// get Fetches a key from the cache if exists.
func (cmd Command) get() bool {
	// check correct no of arguments
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	log.Println("Handle GET")
	val, ok := cache.Get(cmd.Args[1])
	if !ok {
		cmd.Conn.Write([]uint8("$-1\r\n"))
		return true
	}

	// handle quoted strings
	if strings.HasPrefix(val, "\"") {
		val, _ = strconv.Unquote(val)
	}
	// length for debugging
	log.Println("Response length", len(val))
	// send RESP head (eg: $5\r\n)
	cmd.Conn.Write([]uint8(fmt.Sprintf("$%d\r\n", len(val))))
	// convert data to bytes, add newline and send result
	cmd.Conn.Write(append([]uint8(val), []uint8("\r\n")...))

	return true
}
