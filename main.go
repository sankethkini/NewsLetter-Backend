package main

import (
	"context"

	"github.com/sankethkini/NewsLetter-Backend/cmd/app"
)

func main() {
	ctx := context.Background()
	app.Start(ctx)
}
