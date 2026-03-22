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
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	log.Println("Handle GET")
	val, _ := cache.Get(cmd.Args[1])
	if val != nil {
		res, _ := val.(string)
		if strings.HasPrefix(res, "\"") {
			res, _ = strconv.Unquote(res)
		}
		log.Println("Response length", len(res))
		cmd.Conn.Write([]uint8(fmt.Sprintf("$%d\r\n", len(res))))
		cmd.Conn.Write(append([]uint8(res), []uint8("\r\n")...))
	} else {
		cmd.Conn.Write([]uint8("$-1\r\n"))
	}
	return true
}
