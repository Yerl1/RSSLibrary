package client

import (
	"fmt"
	"net"
)

func RunClient(args []string) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	Parser(args, conn)
}

func Parser(args []string, conn net.Conn) {
	switch {
	case args[0] == "fetch":
		Fetch(conn)
	default:
		break
	}
}
