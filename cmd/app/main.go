package main

import (
	"context"
	"log"
	"os"

	"rsslibrary/internal/app"
	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/client"
	"rsslibrary/internal/config"
	"rsslibrary/pkg/loadenv"
)

func main() {
	ctx := context.Background()
	if len(os.Args) == 2 && os.Args[1] == "start_server" {
		loadenv.LoadEnv("./.env")
		cfg := config.Load()
		db, err := repository.ConnectDB(ctx, cfg.Database)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		log.Println("Connected to DB:", db.Stats())
		app.RunApp()
	} else if len(os.Args) > 1 {
		client.RunClient(os.Args[1:])
	}
}
