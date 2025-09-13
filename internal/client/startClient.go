package client

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func RunClient(args []string) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	req := strings.Join(args, " ")
	SendRequest(req, conn)
}
func SendRequest(req string, conn net.Conn) {
	data := []byte(req + "\n")
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
