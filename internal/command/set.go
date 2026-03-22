package command

import (
	"fmt"
	"log"
	"miniredis/internal/cache"
	"strconv"
	"strings"
	"time"
)

// set Stores a key and value on the Cache. Optionally sets expiration on the key.
func (cmd Command) set() bool {
	if len(cmd.Args) < 3 || len(cmd.Args) > 6 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	log.Println("Handle SET")
	log.Println("Value length", len(cmd.Args[2]))
	if len(cmd.Args) > 3 {
		pos := 3
		option := strings.ToUpper(cmd.Args[pos])
		switch option {
		case "NX":
			log.Println("Handle NX")
			if _, ok := cache.Get(cmd.Args[1]); ok {
				cmd.Conn.Write([]uint8("$-1\r\n"))
				return true
			}
			pos++
		case "XX":
			log.Println("Handle XX")
			if _, ok := cache.Get(cmd.Args[1]); !ok {
				cmd.Conn.Write([]uint8("$-1\r\n"))
				return true
			}
			pos++
		}
		if len(cmd.Args) > pos {
			if err := cmd.setExpiration(pos); err != nil {
				cmd.Conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
				return true
			}
		}
	}
	cache.Set(cmd.Args[1], cmd.Args[2])
	cmd.Conn.Write([]uint8("+OK\r\n"))
	return true
}

// setExpiration Handles expiration when passed as part of the 'set' command.
func (cmd Command) setExpiration(pos int) error {
	option := strings.ToUpper(cmd.Args[pos])
	value, _ := strconv.Atoi(cmd.Args[pos+1])
	var duration time.Duration
	switch option {
	case "EX":
		duration = time.Second * time.Duration(value)
	case "PX":
		duration = time.Millisecond * time.Duration(value)
	default:
		return fmt.Errorf("expiration option is not valid")
	}
	go func() {
		log.Printf("Handling '%s', sleeping for %v\n", option, duration)
		time.Sleep(duration)
		cache.Delete(cmd.Args[1])
	}()
	return nil
}
