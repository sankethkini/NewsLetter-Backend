package main

import (
	"context"

	"github.com/sankethkini/NewsLetter-Backend/cmd/app"
)

func main() {
	ctx := context.Background()
	app.Start(ctx)
	// s := email.NewMail([]string{"sankethkini@gmail.com"}, "some", "some")
	// cf, _ := config.LoadConfig()

	// m := email.NewEmailServer(config.LoadEmailConfig(cf))

	// err := m.SendEmail(s)
	// if err != nil {
	// 	panic(err)
	// }
}
