package client

import (
	"fmt"
	"net"
)

func Fetch(conn net.Conn) {
	data := []byte("Fetch")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
