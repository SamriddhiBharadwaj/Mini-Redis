package command

import (
	"fmt"
	"miniredis/internal/cache"
)

// del Deletes a key from the cache.
func (cmd *Command) del() bool {
	count := 0
	for _, k := range cmd.Args[1:] {
		if _, ok := cache.LoadAndDelete(k); ok {
			count++
		}
	}
	cmd.Conn.Write([]uint8(fmt.Sprintf(":%d\r\n", count)))
	return true
}
