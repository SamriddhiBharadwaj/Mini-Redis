package command

import (
	"fmt"
	"miniredis/internal/cache"
	"strconv"
)

func (cmd Command) decr() bool {
	// check len of arguments
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}

	key := cmd.Args[1]

	val, ok := cache.Get(key)

	var num int

	// if key doesnt exist
	if !ok {
		num = 0
	} else {
		var err error
		// check if string is int, else error
		num, err = strconv.Atoi(val)
		if err != nil {
			cmd.Conn.Write([]uint8("-ERR value is not an integer or out of bounds\r\n"))
			return true
		}
	}
	// decrement num
	num--

	// set new values to map
	cache.Set(key, strconv.Itoa(num))

	cmd.Conn.Write([]uint8(fmt.Sprintf(":%d\r\n", num)))

	return true
}
