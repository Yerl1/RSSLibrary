package main

import (
	"context"

	"rsslibrary/internal/app"
)

func main() {
	ctx := context.Background()
	app.RunApp(ctx)
}
