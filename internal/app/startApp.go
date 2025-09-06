package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"

	"rsslibrary/internal/app/handlers"
	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service"
	"rsslibrary/internal/config"
	"rsslibrary/pkg/loadenv"
)

func RunApp(ctx context.Context) {
	runtime.GOMAXPROCS(3)
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()
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

func handleClient(conn net.Conn, ctx context.Context, handler handlers.Handler) {
	defer conn.Close()
	for {
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}
		source := string(input[0:n])
		if source == "fetch" {
			go handler.Fetch(ctx, conn)
		}
		if len(source) >= 13 && source[:12] == "set-interval" {
			go handler.SetInterval(source[13:], ctx, conn)
		}

		if len(source) >= 13 && source[:11] == "set-workers" {
			go handler.SetWorkers(source[12:], ctx, conn)
		}
	}
}
