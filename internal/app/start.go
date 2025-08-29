package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"rsslibrary/internal/app/handlers"
	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service"
	"rsslibrary/internal/config"
	"rsslibrary/pkg/loadenv"
)

func RunApp(ctx context.Context) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server is running on port 8080")
	loadenv.LoadEnv("./.env")
	cfg := config.Load()
	db, err := repository.ConnectDB(ctx, cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("Connected to DB:", db.Stats())
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handlers.NewRequestHandler(service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleClient(conn, ctx, handler)
	}
}

func handleClient(conn net.Conn, ctx context.Context, handler *handlers.RequestHandler) {
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
	fmt.Println(req.String())
	// Processing the request

}
