package main

import (
	"context"
	"fmt"

	subscriptionpb "github.com/sankethkini/NewsLetter-Backend/proto/subscriptionpb/v1"
	userpb "github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"google.golang.org/grpc"
)

// nolint: staticcheck
func usr() {
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

// nolint: staticcheck
func main() {
	serverAddress := "localhost:9000"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		panic(e)
	}
	defer conn.Close()
	cl := subscriptionpb.NewSubscriptionServiceClient(conn)
	data := subscriptionpb.CreateSchemeRequest{
		Name:  "some",
		Price: 500,
		Days:  28,
	}
	ctx := context.Background()
	resp, err := cl.CreateScheme(ctx, &data)
	fmt.Println(resp, err)
}
