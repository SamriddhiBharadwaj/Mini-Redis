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
	value, err := strconv.Atoi(cmd.Args[pos+1])
	if err != nil {
		return fmt.Errorf("invalid expiration value")
	}
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
	log.Printf("Handling '%s'", option)
	cache.SetWithExpiration(cmd.Args[1], cmd.Args[2], duration)
	return nil
}
