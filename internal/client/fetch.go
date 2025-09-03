package client

import (
	"fmt"
	"net"
	"time"
)

func Fetch(conn net.Conn) {
	data := []byte("fetch")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ResponseReader(conn))
}

func ResponseReader(conn net.Conn) string {
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println(err.Error())
		return "error occured"
	}
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 700))
	return string(buff[0:n])
}
