package main

import (
	"fmt"
	"os"
	"rsslibrary/internal/client"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Enter a command to start a program\n for list of commads write\n \t ./rsshub --help")
		return
	}
	client.RunClient(args[1:])
}
