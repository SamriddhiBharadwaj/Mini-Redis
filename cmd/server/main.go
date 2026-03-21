package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sync map over normal map due to concurrent access via goroutines (can also use map + mutex)
var cache sync.Map

func main(){
	server.Start(":6380")
}
