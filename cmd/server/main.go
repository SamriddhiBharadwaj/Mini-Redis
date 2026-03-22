package main

import (
	"miniredis/internal/server"
)

func main() {
	server.Start(":6380")
}
