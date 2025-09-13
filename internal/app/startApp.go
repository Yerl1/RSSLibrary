package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"rsslibrary/internal/app/handlers"
	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service"
	"rsslibrary/internal/config"
	"rsslibrary/pkg/loadenv"
)

func RunApp(ctx context.Context) {
	// runtime.GOMAXPROCS(3)

	// Config
	if err := loadenv.LoadEnv("./.env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	cfg := config.Load()

	// DB
	db, err := repository.ConnectDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	// Services
	repo := repository.NewRepository(db)
	srv := service.NewService(repo)
	handler := handlers.NewRequestHandler(srv)

	// TCP
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	log.Println("Server started at localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handleClient(ctx, conn, handler)
	}
}

func handleClient(ctx context.Context, conn net.Conn, handler handlers.Handler) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v", r)
		}
		conn.Close()
	}()

	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}

		args := strings.Split(line, " ")
		cmd := args[0]
		params := args[1:]

		switch cmd {
		case "fetch":
			go handler.Fetch(ctx, conn)
		case "add":
			name, url := parseFlag(params, "--name"), parseFlag(params, "--url")
			if name == "" || url == "" {
				writeSafe(conn, "Usage: rsshub add --name <name> --url <url>\n")
				continue
			}
			go handler.AddFeed(ctx, name, url, conn)

		case "set-interval":
			if len(params) != 1 {
				writeSafe(conn, "Usage: rsshub set-interval <duration>\n")
				continue
			}
			go handler.SetInterval(ctx, params[0], conn)

		case "set-workers":
			if len(params) != 1 {
				writeSafe(conn, "Usage: rsshub set-workers <number>\n")
				continue
			}
			go handler.SetWorkers(ctx, params[0], conn)

		case "list":
			stringNum := parseFlag(params, "--num")
			num, err := strconv.Atoi(stringNum)
			if err != nil {
				conn.Write([]byte("User input error: " + stringNum + " is not a number"))
				return
			}
			go handler.ListFeeds(ctx, num, conn)

		case "delete":
			name := parseFlag(params, "--name")
			if name == "" {
				writeSafe(conn, "Usage: rsshub delete --name <feed-name>\n")
				continue
			}
			go handler.DeleteFeed(ctx, name, conn)

		case "articles":
			feed := parseFlag(params, "--feed-name")
			if feed == "" {
				writeSafe(conn, "Usage: rsshub articles --feed-name <name> [--num N]\n")
				continue
			}
			stringNum := parseFlag(params, "--num")
			num, err := strconv.Atoi(stringNum)
			if err != nil {
				conn.Write([]byte("User input error: " + stringNum + " is not a number"))
				return
			}
			go handler.ListArticles(ctx, feed, num, conn)

		case "--help", "help":
			writeSafe(conn, cliHelp())

		default:
			writeSafe(conn, fmt.Sprintf("Unknown command: %s\nType rsshub --help for usage.\n", cmd))
		}
	}

	if err := reader.Err(); err != nil {
		log.Printf("client read error: %v", err)
	}
}

func writeSafe(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("failed to write to client: %v", err)
	}
}

func parseFlag(args []string, key string) string {
	for i, arg := range args {
		if arg == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func cliHelp() string {
	return `Usage: rsshub <command> [options]

Commands:
  fetch                          Start background RSS fetcher
  add --name <n> --url <u>       Add new RSS feed
  set-interval <dur>             Set fetch interval (e.g. 2m, 5m)
  set-workers <n>                Set number of workers
  list [--num N]                 List RSS feeds
  delete --name <n>              Delete feed by name
  articles --feed-name <n> [--num N]  Show latest articles
  --help                         Show this help

Examples:
  rsshub fetch
  rsshub add --name tech --url https://techcrunch.com/feed/
  rsshub set-interval 2m
  rsshub set-workers 5
  rsshub list --num 5
  rsshub delete --name tech
  rsshub articles --feed-name tech --num 3
`
}
