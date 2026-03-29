package command

import (
	"fmt"
	"log"
	"miniredis/internal/cache"
	"strconv"
	"strings"
	"time"
)

// setExpiration Handles expiration when passed as part of the 'set' command.
func (cmd Command) setExpiration(pos int) error {
	option := strings.ToUpper(cmd.Args[pos])
	value, _ := strconv.Atoi(cmd.Args[pos+1])
	// create a time object
	var duration time.Duration
	switch option {
	// store expiry in terms of seconds
	case "EX":
		duration = time.Second * time.Duration(value)
	// store expiry in terms of milliseconds
	case "PX":
		duration = time.Millisecond * time.Duration(value)
	default:
		return fmt.Errorf("expiration option is not valid")
	}
	go func() {
		// continuously poll via subroutine
		log.Printf("Handling '%s', sleeping for %v\n", option, duration)
		time.Sleep(duration)
		cache.Delete(cmd.Args[1])
	}()
	return nil
}
