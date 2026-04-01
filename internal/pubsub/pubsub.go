package pubsub

import (
	"fmt"
	"net"
	"sync"
)

// use an array with string as key and stores conn of each subscriber
var channels = make(map[string][]net.Conn)
var mu sync.RWMutex

func Subscribe(channel string, conn net.Conn) {
	// obtain lock
	mu.Lock()
	// append subscriber
	channels[channel] = append(channels[channel], conn)
	mu.Unlock()
	// print to user
	_, err := conn.Write([]uint8(fmt.Sprintf(
		"*3\r\n$9\r\nsubscribe\r\n$%d\r\n%s\r\n:1\r\n",
		len(channel), channel,
	)))
	if err != nil {
		return
	}
	// block forever
	select {}
}

func Publish(channel, message string) int {
	mu.RLock()
	// find subscribers of specific channel
	subs, ok := channels[channel]
	mu.RUnlock()
	if !ok {
		return 0
	}

	count := 0
	// loop over subscribers, and send message to each
	for _, conn := range subs {
		msg := fmt.Sprintf(
			"*3\r\n$7\r\nmessage\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
			len(channel), channel,
			len(message), message,
		)

		_, err := conn.Write([]uint8(msg))

		// if connection has already been closed, handle deletion
		if err != nil {
			mu.Lock()
			newSubs := []net.Conn{}

			// copy ever subscriber except current(inactive) one
			for _, c := range channels[channel] {
				if c != conn {
					newSubs = append(newSubs, c)
				}
			}

			channels[channel] = newSubs
			mu.Unlock()
			continue
		}

		count++
	}
	// return number of subscribers msg is sent to
	return count
}
