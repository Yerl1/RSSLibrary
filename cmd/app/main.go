package main

import (
	"os"

	"rsslibrary/internal/app"
	"rsslibrary/internal/client"
)

func main() {
	// ctx := context.Background()
	// loadenv.LoadEnv("./.env")
	// cfg := config.Load()
	// db, err := repository.ConnectDB(ctx, cfg.Database)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// log.Println("Connected to DB:", db.Stats())

	if len(os.Args) == 2 && os.Args[1] == "start_server" {
		app.RunApp()
	} else if len(os.Args) > 1 {
		client.RunClient(os.Args[1:])
	}
}
