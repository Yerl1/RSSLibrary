package app

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func RunApp() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server is running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Parsing client request
	buffer := make([]byte, 1)
	var req strings.Builder
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}
		req.WriteByte(buffer[0])
	}

	// Processing the request
}
