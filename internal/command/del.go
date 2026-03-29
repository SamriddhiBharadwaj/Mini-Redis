package command

import (
	"fmt"
	"miniredis/internal/cache"
)

// del Deletes a key from the cache.
func (cmd *Command) del() bool {
	// count no. of deleted keys
	count := 0
	// loop over index and key values in slice (ignoring initial "DEL")
	for _, k := range cmd.Args[1:] {
		if _, ok := cache.LoadAndDelete(k); ok {
			count++
		}
	}
	// send response to client
	cmd.Conn.Write([]uint8(fmt.Sprintf(":%d\r\n", count)))
	return true // ensures session is not closed
}
