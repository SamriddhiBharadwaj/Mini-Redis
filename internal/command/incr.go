package command

import (
	"fmt"
	"miniredis/internal/cache"
	"strconv"
)

func (cmd Command) incr() bool {
	if len(cmd.Args) != 2 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "'command\r\n"))
		return true
	}

	key := cmd.Args[1]

	val, ok := cache.Get(key)

	var num int

	// key doesnt exist case
	if !ok {
		num = 0
	} else {
		// check if string is int, else error
		var err error
		num, err = strconv.Atoi(val)
		if err != nil {
			cmd.Conn.Write([]uint8("-ERR value is not an integer or out of range\r\n"))
			return true
		}
	}
	// increment num
	num++

	// store back in cache
	cache.Set(key, strconv.Itoa(num))

	cmd.Conn.Write(([]uint8(fmt.Sprintf(":%d\r\n", num))))
	return true
}
