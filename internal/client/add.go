package client

import (
	"fmt"
	"net"
)

func AddFeed(message string, conn net.Conn) {
	data := []byte(message)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ResponseReader(conn)
}
