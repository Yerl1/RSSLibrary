package client

import (
	"fmt"
	"net"
)

func SetWorkers(number string, conn net.Conn) {
	data := []byte("set-workers " + number)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ResponseReader(conn))
}
