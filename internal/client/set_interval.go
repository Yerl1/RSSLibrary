package client

import (
	"fmt"
	"net"
)

func SetInteval(interval string, conn net.Conn) {
	data := []byte("set-interval " + interval)
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ResponseReader(conn))
}
