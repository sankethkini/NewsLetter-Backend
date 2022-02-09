package app

import (
	"context"
	"fmt"
	"net"

	"github.com/sankethkini/NewsLetter-Backend/internal/service/user"
	"github.com/sankethkini/NewsLetter-Backend/proto/userpb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	transport "github.com/sankethkini/NewsLetter-Backend/internal/transport/user"
)

func Start(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 9000))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	fmt.Printf("listening: %v", err)
	svc := user.NewUserService()
	grpcServer := transport.NewGrpcServer(ctx, svc)
	srv := grpc.NewServer()
	userpb.RegisterUserServiceServer(srv, grpcServer)
	reflection.Register(srv)
	err = srv.Serve(lis)
	if err != nil {
		fmt.Println(err)
	}

}
