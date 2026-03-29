package command

import (
	"log"
	"miniredis/internal/cache"
	"strings"
)

// set Stores a key and value on the Cache. Optionally sets expiration on the key.
func (cmd Command) set() bool {
	// error check
	if len(cmd.Args) < 3 || len(cmd.Args) > 6 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	// logs for debugging
	log.Println("Handle SET")
	log.Println("Value length", len(cmd.Args[2]))
	// check for special flags
	if len(cmd.Args) > 3 {
		pos := 3
		option := strings.ToUpper(cmd.Args[pos])
		switch option {
		// only set the key if it does not already exist
		case "NX":
			log.Println("Handle NX")
			// if key exists, donot set
			if _, ok := cache.Get(cmd.Args[1]); ok {
				cmd.Conn.Write([]uint8("$-1\r\n"))
				return true
			}
			pos++
		// only set the key if it already exist
		case "XX":
			log.Println("Handle XX")
			// if key doesnt exist, donot set
			if _, ok := cache.Get(cmd.Args[1]); !ok {
				cmd.Conn.Write([]uint8("$-1\r\n"))
				return true
			}
			pos++
		}
		// handle expiration
		if len(cmd.Args) > pos {
			// handle error
			if err := cmd.setExpiration(pos); err != nil {
				cmd.Conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
				return true
			}
		}
	}
	// set to map
	cache.Set(cmd.Args[1], cmd.Args[2])
	cmd.Conn.Write([]uint8("+OK\r\n"))
	return true
}
