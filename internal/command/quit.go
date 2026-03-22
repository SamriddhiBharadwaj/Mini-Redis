package command

import "log"

func (cmd *Command) quit() bool {
	if len(cmd.Args) != 1 {
		cmd.Conn.Write([]uint8("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n"))
		return true
	}
	log.Println("Handle QUIT")
	cmd.Conn.Write([]uint8("+OK\r\n"))
	return false
}
