package main

import (
	"context"
	"fmt"

	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"google.golang.org/grpc"
)

func main() {
	serverAddress := "localhost:9000"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		panic(e)
	}
	defer conn.Close()
	cl := userpb.NewUserServiceClient(conn)
	data := userpb.CreateUserRequest{
		User: &userpb.User{
			Email:    "somemail@some.com",
			Password: "some password",
			Name:     "sank",
		},
	}
	data1 := userpb.ValidateUserRequest{
		Email:    "somemail@some.com",
		Password: "some password",
	}
	ctx := context.Background()
	resp, err := cl.CreateUser(ctx, &data)
	fmt.Println(resp, err)
	resp1, err := cl.ValidateUser(ctx, &data1)
	fmt.Println(resp1, err)
}
