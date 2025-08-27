package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"rsslibrary/internal/config"
	"rsslibrary/internal/repository"
	"rsslibrary/pkg/loadenv"
)

func main() {
	ctx := context.Background()
	loadenv.LoadEnv("./.env")
	cfg := config.Load()
	db, err := repository.ConnectDB(ctx, cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("Connected to DB:", db.Stats())

	for {
		if len(os.Args) > 1 {
			fmt.Println(os.Args)
		}
	}
}
