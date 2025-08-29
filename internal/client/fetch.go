package client

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func Fetch(conn net.Conn) {
	data := []byte("fetch")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Println(ResponseReader(conn))
}

func ResponseReader(conn net.Conn) string {
	var resp strings.Builder
	buffer := make([]byte, 1)
	fmt.Println("Waiting 1")
	for {
		_, err := conn.Read(buffer)
		fmt.Println("Waiting 2")
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return ""
		}
		resp.WriteByte(buffer[0])
	}

	return resp.String()
}
