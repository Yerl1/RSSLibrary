package main

import (
	"context"
	"os"

	"rsslibrary/internal/app"
	"rsslibrary/internal/client"
)

func main() {
	ctx := context.Background()
	if len(os.Args) == 2 && os.Args[1] == "start_server" {

		app.RunApp(ctx)
	} else if len(os.Args) > 1 {
		client.RunClient(os.Args[1:])
	}
}
