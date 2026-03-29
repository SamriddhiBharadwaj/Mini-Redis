package main

import (
	"miniredis/internal/cache"
	"miniredis/internal/server"
)

func main() {
	cache.ActiveExpiration()
	server.Start(":6380")
}
