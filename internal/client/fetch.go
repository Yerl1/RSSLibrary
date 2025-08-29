package client

import (
	"fmt"
	"net"
)

func Fetch(conn net.Conn) {
	data := []byte("fetch")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
